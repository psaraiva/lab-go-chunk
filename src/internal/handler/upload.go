package handler

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"lab/src/internal/entity"
	"lab/src/logger"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (ac *Action) FeatureUpload() error {
	logger.GetLogActivity().WriteLog("Gerando HASH do arquivo...")
	err := ac.GenerateHashFile()
	if err != nil {
		ac.LogErrorWrite(
			ACTION_UPLOAD,
			fmt.Sprintf(
				"fail to generate hash file: %s",
				ac.FileTarget))

		if err == errorFileNotFound {
			return fmt.Errorf("arquivo não encontrado: %s", ac.FileTarget)
		}

		return errorActionDefault
	}

	logger.GetLogActivity().WriteLog("Validando arquivo...")
	flag, err := ac.isNewFile()
	if err != nil {
		ac.LogErrorWrite(
			ACTION_UPLOAD,
			fmt.Sprintf(
				"fail to verify is new file: %s: %s",
				ac.FileTarget,
				err.Error()))

		return errorActionDefault
	}

	if flag {
		ac.LogErrorWrite(
			ACTION_UPLOAD,
			fmt.Sprintf(
				"the file already exists: %s",
				ac.FileTarget))

		return fmt.Errorf("o arquivo já existe: %s", ac.FileTarget)
	}

	file := entity.File{
		Hash: ac.Hash,
		Name: filepath.Base(ac.FileTarget),
	}

	_, err = repositoryFile.Create(file)
	if err != nil {
		ac.LogErrorWrite(
			ACTION_UPLOAD,
			fmt.Sprintf(
				"fail to create file by repository: %s",
				err.Error()))
		return errorActionDefault
	}

	logger.GetLogActivity().WriteLog("Copiando arquivo original para área temporária...")
	err = ac.SendFileToTemporaryArea()
	if err != nil {
		ac.LogErrorWrite(
			ACTION_UPLOAD,
			fmt.Sprintf(
				"fail to move file for temporary area: %s",
				err.Error()))
		return errorActionDefault
	}

	logger.GetLogActivity().WriteLog("Gerando chunks do arquivo original em memória...")
	chunk, err := ac.GenerateChunkByHashFile(ac.Hash)
	if err != nil {
		ac.LogErrorWrite(
			ACTION_UPLOAD,
			fmt.Sprintf(
				"fail to generate chunck by hash file: %s: %s",
				ac.Hash,
				err.Error()))
		return errorActionDefault
	}

	logger.GetLogActivity().WriteLog("Salvando coleção de chunk...")
	_, err = repositoryChunk.Create(chunk)
	if err != nil {
		ac.LogErrorWrite(
			ACTION_UPLOAD,
			fmt.Sprintf(
				"fail to create chunk by repository: %s",
				err.Error()))
		return errorActionDefault
	}

	logger.GetLogActivity().WriteLog("Gerando arquivos chunks .bin em storage...")
	err = ac.GenerateChunksToStorage(chunk)
	if err != nil {
		ac.LogErrorWrite(
			ACTION_UPLOAD,
			fmt.Sprintf(
				"fail to generate chunks to storage: %s",
				err.Error()))
		return errorActionDefault
	}

	logger.GetLogActivity().WriteLog("Removendo arquivo original da área temporária...")
	err = serviceTemporaryArea.RemoveFile(ac.Hash)
	if err != nil {
		ac.LogErrorWrite(
			ACTION_UPLOAD,
			fmt.Sprintf(
				"fail to remove file for temporary area: %s",
				err.Error()))
		return errorActionDefault
	}

	return nil
}

func (ac Action) GenerateChunksToStorage(chunk entity.Chunk) error {
	file, err := serviceTemporaryArea.GetFile(chunk.HashOriginalFile)
	if err != nil {
		return err
	}
	defer file.Close()

	index := -1
	buf := make([]byte, chunk.Size)
	for {
		index++
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		err = serviceStorage.CreateFile(buf[:n], chunk.HashList[index]+".bin")
		if err != nil {
			return err
		}
	}

	return nil
}

func (ac Action) GenerateChunkByHashFile(hashFile string) (entity.Chunk, error) {
	chunk := entity.Chunk{}
	file, err := serviceTemporaryArea.GetFile(hashFile)
	if err != nil {
		return chunk, err
	}
	defer file.Close()

	chunkSizeStr := os.Getenv("CHUNK_SIZE")
	chunkSize, err := strconv.Atoi(chunkSizeStr)
	if err != nil {
		return chunk, fmt.Errorf("configuration CHUNK_SIZE")
	}

	var chunks []string
	buf := make([]byte, chunkSize)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return chunk, err
		}

		if n == 0 {
			break
		}

		chunkHash := md5.Sum(buf[:n])
		chunkHashString := hex.EncodeToString(chunkHash[:])
		chunks = append(chunks, chunkHashString)
	}

	chunk.HashOriginalFile = hashFile
	chunk.HashList = chunks
	chunk.Size = chunkSize
	return chunk, nil
}

func (ac *Action) GenerateHashFile() error {
	file, err := os.Open(ac.FileTarget)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return errorFileNotFound
		}

		return err
	}
	defer file.Close()

	hash, err := entity.File{}.GenerateHashByOsFile(file)
	if err != nil {
		return err
	}

	ac.Hash = hash
	return nil
}

func (ac *Action) SendFileToTemporaryArea() error {
	source, err := os.Open(ac.FileTarget)
	if err != nil {
		return err
	}
	defer source.Close()
	return serviceTemporaryArea.CreateFileByFileSource(ac.Hash, source)
}

func (ac Action) isNewFile() (bool, error) {
	return repositoryFile.IsExistsByHash(ac.Hash)
}

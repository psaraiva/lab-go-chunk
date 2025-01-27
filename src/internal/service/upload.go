package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"lab/src/logger"
	"lab/src/model"
	"os"
	"path/filepath"
	"strconv"
)

func (ac *Action) FeatureUpload() error {
	logger.GetLogActivity().WriteLog("Gerando HASH do arquivo...")
	err := ac.GenerateHashFile()
	if err != nil {
		return err
	}

	logger.GetLogActivity().WriteLog("Validando arquivo...")
	err = ac.isNewFile()
	if err != nil {
		return err
	}

	file := model.File{
		Hash: ac.Hash,
		Name: filepath.Base(ac.FileTarget),
	}

	_, err = repositoryFile.Create(file)
	if err != nil {
		logger.GetLogError().WriteLog(err.Error())
		return fmt.Errorf("erro ao processar o HASH do arquivo")
	}

	logger.GetLogActivity().WriteLog("Copiando arquivo original para área temporária...")
	err = ac.SendFileToTmp()
	if err != nil {
		logger.GetLogError().WriteLog(err.Error())
		return err
	}

	logger.GetLogActivity().WriteLog("Gerando chunks do arquivo original em memória...")
	chunk, err := ac.GenerateChunkByHashFile(ac.Hash)
	if err != nil {
		logger.GetLogError().WriteLog(err.Error())
	}

	logger.GetLogActivity().WriteLog("Salvando coleção de chunk...")
	_, err = repositoryChunk.Create(chunk)
	if err != nil {
		logger.GetLogError().WriteLog(err.Error())
		return err
	}

	logger.GetLogActivity().WriteLog("Gerando arquivos chunks .bin em storage...")
	err = ac.GenerateChunksToStorage(chunk)
	if err != nil {
		logger.GetLogError().WriteLog(err.Error())
		return err
	}

	logger.GetLogActivity().WriteLog("Removendo arquivo original da área temporária...")
	err = serviceTemporaryArea.RemoveFile(ac.Hash)
	if err != nil {
		logger.GetLogError().WriteLog(err.Error())
		return err
	}

	return nil
}

func (ac Action) GenerateChunksToStorage(chunk model.Chunk) error {
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

func (ac Action) GenerateChunkByHashFile(hashFile string) (model.Chunk, error) {
	chunk := model.Chunk{}
	file, err := serviceTemporaryArea.GetFile(hashFile)
	if err != nil {
		return chunk, err
	}
	defer file.Close()

	chunkSizeStr := os.Getenv("CHUNK_SIZE")
	chunkSize, err := strconv.Atoi(chunkSizeStr)
	if err != nil {
		return chunk, fmt.Errorf("falha na configuração de: Chunk Size")
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
		return err
	}
	defer file.Close()

	hash, err := model.File{}.GenerateHashByOsFile(file)
	if err != nil {
		return err
	}

	ac.Hash = hash
	return nil
}

func (ac *Action) SendFileToTmp() error {
	source, err := os.Open(ac.FileTarget)
	if err != nil {
		return err
	}
	defer source.Close()
	return serviceTemporaryArea.CreateFileByFileSource(ac.Hash, source)
}

func (ac Action) isNewFile() error {
	resp, err := repositoryFile.IsExistsByHashFile(ac.Hash)
	if resp {
		return fmt.Errorf("arquivo já existe: %s", ac.FileTarget)
	}

	return err
}

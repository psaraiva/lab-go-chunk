package actions

import (
	"crypto/md5"
	"encoding/binary"
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

	err = repositoryFile.Create(file)
	if err != nil {
		logger.GetLogError().WriteLog(err.Error())
		return fmt.Errorf("erro ao processar o HASH do arquivo")
	}

	logger.GetLogActivity().WriteLog("Copiando arquivo para pasta temporária...")
	err = ac.SendFileToTmp()
	if err != nil {
		logger.GetLogError().WriteLog(err.Error())
		return err
	}

	logger.GetLogActivity().WriteLog("Gerando chunks do arquivo...")
	chunkItem, err := ac.GenerateChunkByHashFile(ac.Hash)
	if err != nil {
		logger.GetLogError().WriteLog(err.Error())
	}

	logger.GetLogActivity().WriteLog("Salvando chunks do arquivo...")
	err = repositoryChunkItem.Create(chunkItem)
	if err != nil {
		logger.GetLogError().WriteLog(err.Error())
		return err
	}

	err = ac.GenerateChunksFileTmp(chunkItem)
	if err != nil {
		logger.GetLogError().WriteLog(err.Error())
		return err
	}

	logger.GetLogActivity().WriteLog("Removendo arquivo temporário...")
	err = ac.removeFileToTmp()
	if err != nil {
		logger.GetLogError().WriteLog(err.Error())
		return err
	}

	return nil
}

func (ac Action) GenerateChunksFileTmp(chunkItem model.ChunkItem) error {
	file, err := os.Open(os.Getenv("FOLDER_TMP") + "/" + chunkItem.HashFile)
	if err != nil {
		return err
	}
	defer file.Close()

	index := -1
	buf := make([]byte, chunkItem.Size)
	for {
		index++
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		err = ac.saveChunkBin(buf[:n], chunkItem.HashList[index])
		if err != nil {
			return err
		}
	}

	return nil
}

func (ac Action) GenerateChunkByHashFile(hashFile string) (model.ChunkItem, error) {
	file, err := os.Open(os.Getenv("FOLDER_TMP") + "/" + hashFile)
	if err != nil {
		return model.ChunkItem{}, err
	}
	defer file.Close()

	chunkSizeStr := os.Getenv("CHUNK_SIZE")
	chunkSize, err := strconv.Atoi(chunkSizeStr)
	if err != nil {
		return model.ChunkItem{}, fmt.Errorf("falha na configuração de: Chunk Size")
	}

	var chunks []string
	buf := make([]byte, chunkSize)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return model.ChunkItem{}, err
		}

		if n == 0 {
			break
		}

		chunkHash := md5.Sum(buf[:n])
		chunkHashString := hex.EncodeToString(chunkHash[:])
		chunks = append(chunks, chunkHashString)
	}

	chunkModel := model.ChunkItem{}
	chunkModel.HashFile = hashFile
	chunkModel.HashList = chunks
	chunkModel.Size = chunkSize
	return chunkModel, nil
}

func (ac *Action) GenerateHashFile() error {
	file, err := os.Open(ac.FileTarget)
	if err != nil {
		return err
	}
	defer file.Close()

	hashString, err := model.File{}.GenerateHashByOsFile(file)
	if err != nil {
		return err
	}

	ac.Hash = hashString
	return nil
}

func (ac *Action) SendFileToTmp() error {
	sourceFile, err := os.Open(ac.FileTarget)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(os.Getenv("FOLDER_TMP") + "/" + ac.Hash)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func (ac Action) isNewFile() error {
	resp, err := repositoryFile.IsExistsByHashFile(ac.Hash)
	if resp {
		return fmt.Errorf("arquivo já existe: %s", ac.FileTarget)
	}

	return err
}

// @todo Passar esse método para responsabilidade de Service TMP
func (ac *Action) saveChunkBin(data []byte, hash string) error {
	file, err := os.Create(fmt.Sprintf("%s/%s.bin", os.Getenv("FOLDER_STORAGE"), hash))
	if err != nil {
		return err
	}
	defer file.Close()

	for _, value := range data {
		err = binary.Write(file, binary.LittleEndian, value)
		if err != nil {
			return err
		}
	}

	return nil
}

// @todo Passar esse método para responsabilidade de Service TMP
func (ac *Action) removeFileToTmp() error {
	return os.Remove(fmt.Sprintf("%s/%s", os.Getenv("FOLDER_TMP"), ac.Hash))
}

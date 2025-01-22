package actions

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"lab/src/logger"
	"os"
	"path/filepath"
)

func (ac *Action) FeatureUpload() error {
	logger.GetLogActivity().WriteLog("Processando HASH do arquivo...")

	err := ac.isNewFile()
	if err != nil {
		return err
	}

	err = ac.addHashToCollection()
	if err != nil {
		logger.GetLogError().WriteLog(err.Error())
		return fmt.Errorf("erro ao processar o HASH do arquivo")
	}

	logger.GetLogActivity().WriteLog("Copiando arquivo para pasta temporária...")
	err = ac.sendFileToTmp()
	if err != nil {
		logger.GetLogError().WriteLog(err.Error())
		return err
	}

	logger.GetLogActivity().WriteLog("Processando partes do arquivo...")
	err = ac.processChunk()
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

func (ac *Action) isNewFile() error {
	err := ac.GenerateHashFile()
	if err != nil {
		return err
	}

	jsonHash, err := os.Open(os.Getenv("JSON_FILE_HASH"))
	if err != nil {
		return err
	}
	defer jsonHash.Close()

	decoder := json.NewDecoder(jsonHash)
	hashList := []hashItem{}
	err = decoder.Decode(&hashList)
	if err != nil {
		return err
	}

	for _, item := range hashList {
		if item.Hash == ac.Hash {
			return fmt.Errorf("arquivo já existe: %s", item.Name)
		}
	}

	return nil
}

func (ac *Action) GenerateHashFile() error {
	file, err := os.Open(ac.FileTarget)
	if err != nil {
		return err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}

	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)
	ac.Hash = hashString

	return nil
}

func (ac *Action) addHashToCollection() error {
	jsonHash, err := os.Open(os.Getenv("JSON_FILE_HASH"))
	if err != nil {
		return err
	}
	defer jsonHash.Close()

	decoder := json.NewDecoder(jsonHash)
	hashList := []hashItem{}
	err = decoder.Decode(&hashList)
	if err != nil {
		return err
	}

	for _, item := range hashList {
		if item.Hash == ac.Hash {
			return fmt.Errorf("arquivo já existe: %s", item.Name)
		}
	}

	hashList = append(hashList, hashItem{Hash: ac.Hash, Name: filepath.Base(ac.FileTarget)})
	updatedJSON, err := json.MarshalIndent(hashList, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(os.Getenv("JSON_FILE_HASH"), updatedJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (ac *Action) sendFileToTmp() error {
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

func (ac *Action) processChunk() error {
	file, err := os.Open(os.Getenv("FOLDER_TMP") + "/" + ac.Hash)
	if err != nil {
		return err
	}
	defer file.Close()

	var chunks []string
	buf := make([]byte, CHUNK_SIZE)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		chunkHash := md5.Sum(buf[:n])
		chunkHashString := hex.EncodeToString(chunkHash[:])
		chunks = append(chunks, chunkHashString)

		err = ac.saveChunkBin(buf[:n], chunkHashString)
		if err != nil {
			return err
		}
	}

	item := chunkItem{}
	item.HashFile = ac.Hash
	item.Chunk = chunks

	jsonChunkFile, err := os.Open(os.Getenv("JSON_FILE_CHUNK"))
	if err != nil {
		return err
	}
	defer jsonChunkFile.Close()

	decoder := json.NewDecoder(jsonChunkFile)
	chunkList := []chunkItem{}
	err = decoder.Decode(&chunkList)
	if err != nil {
		return err
	}

	chunkList = append(chunkList, item)
	updatedJSON, err := json.MarshalIndent(chunkList, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(os.Getenv("JSON_FILE_CHUNK"), updatedJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

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

func (ac *Action) removeFileToTmp() error {
	return os.Remove(fmt.Sprintf("%s/%s", os.Getenv("FOLDER_TMP"), ac.Hash))
}

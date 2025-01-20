package feature

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	ACTION_CLEAR    = "clear"
	ACTION_DOWNLOAD = "download"
	ACTION_REMOVE   = "remove"
	ACTION_UPLOAD   = "upload"
	CHUNK_SIZE      = 1 * 1024 // 1Kb
)

type Action struct {
	Action     string
	FileTarget string
	Hash       string
}

type hashFile struct {
	List []hashFileItem
}

type hashFileItem struct {
	Hash string
	Name string
}

type chunkFile struct {
	List []chunkFileItem
}

type chunkFileItem struct {
	HashFile string
	Chunk    []string
}

func MakeAction() Action {
	return Action{Action: "", FileTarget: "", Hash: ""}
}

func (ac *Action) Execute() error {
	switch ac.Action {
	case ACTION_CLEAR:
		GetLogActivity().WriteLog("Iniciando rotina de limpeza.")
		err := ac.actionClear()
		if err != nil {
			GetLogError().WriteLog(err.Error())
			return err
		}

		GetLogActivity().WriteLog("Rotina de limpeza realizada com sucesso!")
	case ACTION_UPLOAD:
		GetLogActivity().WriteLog("Iniciando rotina de carregamento de arquivo.")
		err := ac.actionUpload()
		if err != nil {
			GetLogError().WriteLog(err.Error())
			return err
		}

		GetLogActivity().WriteLog("Arquivo carregado com sucesso!")
	case ACTION_DOWNLOAD:
		GetLogActivity().WriteLog("Iniciando rotina de descarregamento de arquivo.")
		err := ac.actionDownload()
		if err != nil {
			GetLogError().WriteLog(err.Error())
			return err
		}

		GetLogActivity().WriteLog("Arquivo descarregado com sucesso!")
	case ACTION_REMOVE:
		GetLogActivity().WriteLog("Iniciando rotina de remoção de arquivo.")
		err := ac.actionRemove()
		if err != nil {
			GetLogError().WriteLog(err.Error())
			return err
		}

		GetLogActivity().WriteLog("Arquivo removido com sucesso!")
	}

	return nil
}

func (ac *Action) actionRemove() error {
	GetLogActivity().WriteLog("Carregando Hash do arquivo...")
	hash, err := ac.getHashByFileName(ac.FileTarget)
	if err != nil {
		return err
	}

	GetLogActivity().WriteLog("Carregando lista de chunk do arquivo...")
	list, err := ac.getChunksByHash(hash)
	if err != nil {
		return err
	}

	for _, hashChunk := range list {
		flag, err := ac.isChunkCanBeRemoved(hashChunk)
		if err != nil {
			return err
		}

		if !flag {
			continue
		}

		GetLogActivity().WriteLog("Removendo chunk: " + hashChunk + ".bin")
		err = os.Remove(os.Getenv("FOLDER_STORAGE") + "/" + hashChunk + ".bin")
		if err != nil {
			return err
		}
	}

	GetLogActivity().WriteLog("Removendo registro do arquivo de chunk...")
	err = ac.removeHashToChunkFile(hash)
	if err != nil {
		return err
	}

	GetLogActivity().WriteLog("Removendo registro do arquivo de hash...")
	return ac.removeHashToHashFile(hash)
}

func (ac *Action) removeHashToHashFile(hashString string) error {
	jsonHashFile, err := os.Open(os.Getenv("JSON_FILE_HASH"))
	if err != nil {
		return err
	}
	defer jsonHashFile.Close()

	decoder := json.NewDecoder(jsonHashFile)
	hashFile := hashFile{}
	err = decoder.Decode(&hashFile)
	if err != nil {
		return err
	}

	for index, item := range hashFile.List {
		if item.Hash == hashString {
			hashFile.List = append(hashFile.List[:index], hashFile.List[index+1:]...)
			break
		}
	}

	upJSON, err := json.MarshalIndent(hashFile, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(os.Getenv("JSON_FILE_HASH"), upJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (ac *Action) removeHashToChunkFile(hashFile string) error {
	jsonChunkFile, err := os.Open(os.Getenv("JSON_FILE_CHUNK"))
	if err != nil {
		return err
	}
	defer jsonChunkFile.Close()

	decoder := json.NewDecoder(jsonChunkFile)
	chunkFile := chunkFile{}
	err = decoder.Decode(&chunkFile)
	if err != nil {
		return err
	}

	for index, item := range chunkFile.List {
		if item.HashFile == hashFile {
			chunkFile.List = append(chunkFile.List[:index], chunkFile.List[index+1:]...)
			break
		}
	}

	upJSON, err := json.MarshalIndent(chunkFile, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(os.Getenv("JSON_FILE_CHUNK"), upJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (ac *Action) isChunkCanBeRemoved(chunk string) (bool, error) {
	jsonChunkFile, err := os.Open(os.Getenv("JSON_FILE_CHUNK"))
	if err != nil {
		return false, err
	}
	defer jsonChunkFile.Close()

	decoder := json.NewDecoder(jsonChunkFile)
	chunkFile := chunkFile{}
	err = decoder.Decode(&chunkFile)
	if err != nil {
		return false, err
	}

	countChuks := ac.getCountChunkMap(chunkFile)
	if countChuks[chunk] > 1 {
		return false, nil
	}

	return true, nil
}

func (ac *Action) getCountChunkMap(chunkFile chunkFile) map[string]int {
	chunkCount := make(map[string]int)
	for _, item := range chunkFile.List {
		for _, value := range item.Chunk {
			chunkCount[value]++
		}
	}
	return chunkCount
}

func (ac *Action) actionClear() error {
	err := ac.restoreFileConfigChunk()
	if err != nil {
		return err
	}

	err = ac.restoreFileConfigHash()
	if err != nil {
		return err
	}

	err = ac.cleanStorage()
	if err != nil {
		return err
	}

	err = ac.cleanTmp()
	if err != nil {
		return err
	}

	err = GetLogError().ClearLog()
	if err != nil {
		return err
	}

	err = GetLogActivity().ClearLog()
	if err != nil {
		return err
	}

	return nil
}

func (ac *Action) restoreFileConfigChunk() error {
	err := os.WriteFile(os.Getenv("JSON_FILE_CHUNK"), []byte("{\"List\":[]}"), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (ac *Action) restoreFileConfigHash() error {
	err := os.WriteFile(os.Getenv("JSON_FILE_HASH"), []byte("{\"List\":[]}"), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (ac *Action) cleanTmp() error {
	err := filepath.Walk(os.Getenv("FOLDER_TMP"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		return os.Remove(path)
	})
	return err
}

func (ac *Action) cleanStorage() error {
	err := filepath.Walk(os.Getenv("FOLDER_STORAGE"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		return os.Remove(path)
	})
	return err
}

func (ac *Action) actionUpload() error {
	GetLogActivity().WriteLog("Processando HASH do arquivo...")
	err := ac.processHash()
	if err != nil {
		GetLogError().WriteLog(err.Error())
		return fmt.Errorf("erro ao processar o HASH do arquivo")
	}

	GetLogActivity().WriteLog("Copiando arquivo para pasta temporária...")
	err = ac.sendFileToTmp()
	if err != nil {
		GetLogError().WriteLog(err.Error())
		return err
	}

	GetLogActivity().WriteLog("Processando partes do arquivo...")
	err = ac.processChunk()
	if err != nil {
		GetLogError().WriteLog(err.Error())
		return err
	}

	GetLogActivity().WriteLog("Removendo arquivo temporário...")
	err = ac.removeFileToTmp()
	if err != nil {
		GetLogError().WriteLog(err.Error())
		return err
	}

	return nil
}

func (ac *Action) actionDownload() error {
	hash, err := ac.getHashByFileName(ac.FileTarget)
	if err != nil {
		return err
	}

	list, err := ac.getChunksByHash(hash)
	if err != nil {
		return err
	}

	return ac.generateFileByChunks(list, "")
}

func (ac *Action) generateFileByChunks(chunks []string, targetFolder string) error {
	out, err := os.Create(targetFolder + ac.FileTarget)
	if err != nil {
		return err
	}
	defer out.Close()

	for _, file := range chunks {
		in, err := os.Open(os.Getenv("FOLDER_STORAGE") + "/" + file + ".bin")
		if err != nil {
			return err
		}
		defer in.Close()

		_, err = io.Copy(out, in)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ac *Action) getHashByFileName(fileName string) (string, error) {
	jsonFile, err := os.Open(os.Getenv("JSON_FILE_HASH"))
	if err != nil {
		return "", err
	}
	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)
	hashFile := hashFile{}
	err = decoder.Decode(&hashFile)
	if err != nil {
		return "", err
	}

	for _, item := range hashFile.List {
		if item.Name == fileName {
			return item.Hash, nil
		}
	}

	return "", fmt.Errorf("arquivo não encontrado")
}

func (ac *Action) getChunksByHash(hash string) ([]string, error) {
	jsonFile, err := os.Open(os.Getenv("JSON_FILE_CHUNK"))
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)
	chunkFile := chunkFile{}
	err = decoder.Decode(&chunkFile)
	if err != nil {
		return nil, err
	}

	for _, item := range chunkFile.List {
		if item.HashFile == hash {
			return item.Chunk, nil
		}
	}

	return nil, fmt.Errorf("arquivo não encontrado")
}

func (ac *Action) processHash() error {
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
	return ac.addHashToFile(hashString)
}

func (a *Action) addHashToFile(hash string) error {
	jsonFile, err := os.Open(os.Getenv("JSON_FILE_HASH"))
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)
	hashFile := hashFile{}
	err = decoder.Decode(&hashFile)
	if err != nil {
		return err
	}

	for _, item := range hashFile.List {
		if item.Hash == hash {
			return fmt.Errorf("arquivo já existe")
		}
	}

	hashFile.List = append(hashFile.List, hashFileItem{Hash: hash, Name: filepath.Base(a.FileTarget)})
	updatedJSON, err := json.MarshalIndent(hashFile, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(os.Getenv("JSON_FILE_HASH"), updatedJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (a *Action) sendFileToTmp() error {
	sourceFile, err := os.Open(a.FileTarget)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(os.Getenv("FOLDER_TMP") + "/" + a.Hash)
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

func (a *Action) processChunk() error {
	file, err := os.Open(os.Getenv("FOLDER_TMP") + "/" + a.Hash)
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

		err = a.saveChunkBin(buf[:n], chunkHashString)
		if err != nil {
			return err
		}
	}

	item := chunkFileItem{}
	item.HashFile = a.Hash
	item.Chunk = chunks

	jsonChunkFile, err := os.Open(os.Getenv("JSON_FILE_CHUNK"))
	if err != nil {
		return err
	}
	defer jsonChunkFile.Close()

	decoder := json.NewDecoder(jsonChunkFile)
	chunkFile := chunkFile{}
	err = decoder.Decode(&chunkFile)
	if err != nil {
		return err
	}

	chunkFile.List = append(chunkFile.List, item)
	updatedJSON, err := json.MarshalIndent(chunkFile, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(os.Getenv("JSON_FILE_CHUNK"), updatedJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (a *Action) saveChunkBin(data []byte, hash string) error {
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

func (a *Action) removeFileToTmp() error {
	return os.Remove(fmt.Sprintf("%s/%s", os.Getenv("FOLDER_TMP"), a.Hash))
}

package env

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Env map type
type EnvMap map[string]string

// Map is a struct that contains the map parsed lines from the env file and methods to update the map
type Map struct {
	Map EnvMap
}

// Sets the key and value to the map
func (e *Map) Set(key, val string) {
	e.Map[key] = val
}

// SetMap reads all keys and values from the EnvMap and sets to another EnvMap.Env
func (e *Map) SetMap(target *Map) {
	for key, val := range target.Map {
		e.Set(key, val)
	}
}

// NewEnvMap creates and returns EnvMap, as well as cretaing the map
func NewMap() *Map {
	m := &Map{}

	m.Map = make(EnvMap)
	return m
}

// Adapter is a interface for pulling  secrets from a external secerts storage service (ex. AWS secret manager) and exporting them in your application
type Adapter struct {
	// Pull fucntion will be where secrets will be retrieved and will return a EnvMap
	Pull func() (*Map, error)
}

var (
	envFileNames = []string{".env"}
	requiredKeys []string
	adapters     []*Adapter
)

/* Load scans one or mores that are given and exports the vairbles in the file if they do not exist.
if a file is not provided then the `.env` file in the current working directory will be scaned
instead if one was found.

After that happens load will run the adapters if any were provided then it will run thos adapters
to return a env map that will be exported as well
*/
func Load(filenames ...string) error {
	if len(filenames) == 0 {
		filenames = envFileNames
	}

	// load files
	files, err := loadFiles(false, filenames...)
	if err != nil {
		return err
	}

	globalEnvMap := NewMap()

	// parse files
	for _, content := range files {
		// parse file
		emap := Parse(content)

		globalEnvMap.SetMap(emap)
	}

	if len(adapters) != 0 {
		// run pull secrets from adapters
		for _, adapter := range adapters {

			// pulling secrets
			emap, err := adapter.Pull()
			if err != nil {
				return fmt.Errorf("error occured running adapter: %s", err)
			}

			// set adapters EnvMap to global EnvMap
			globalEnvMap.SetMap(emap)
		}
	}

	// set env map to env
	err = setEnvMap(globalEnvMap)
	if err != nil {
		return err
	}

	return nil
}

/* Load scans one or mores that are given and exports the vairbles in the file if they do not exist.
if a file is not provided then the `.env` file in the current working directory will be scaned
instead if one was found.

After that happens load will run the adapters if any were provided then it will run thos adapters
to return a env map that will be exported as well.

This will error if a required key/s are missing if require keys were provided
*/
func MustLoad(filenames ...string) error {
	err := Load(filenames...)
	if err != nil {
		return err
	}

	// check for missing required keys
	if len(requiredKeys) != 0 {
		var missingKeys []string

		for _, key := range requiredKeys {
			val, ok := os.LookupEnv(key)

			if !ok && val == "" {
				missingKeys = append(missingKeys, key)
			}

		}

		if len(missingKeys) != 0 {
			return fmt.Errorf("Required keys missing or empty: %s", missingKeys)
		}
	}

	return nil
}

// LoadSecrets will run all your adapters and set all the env vars that were fetch then set them to your env in your application
func LoadSecrets() error {
	globalEnvMap := NewMap()

	if len(adapters) != 0 {
		// run pull secrets from adapters
		for _, adapter := range adapters {

			// pulling secrets
			emap, err := adapter.Pull()
			if err != nil {
				return fmt.Errorf("error occured running adapter: %s", err)
			}

			// set adapters EnvMap to global EnvMap
			globalEnvMap.SetMap(emap)
		}
	}

	// set env map to env
	err := setEnvMap(globalEnvMap)
	if err != nil {
		return err
	}

	return nil
}

/* Must LoadSecrets will run all your adapters and set all the env vars that were fetch then set them to your env in your application.
As well as checking for required secrets */
func MustLoadSecrets() error {
	err := LoadSecrets()
	if err != nil {
		return err
	}

	// check for missing required keys
	if len(requiredKeys) != 0 {
		var missingKeys []string

		for _, key := range requiredKeys {
			val, ok := os.LookupEnv(key)

			if !ok && val == "" {
				missingKeys = append(missingKeys, key)
			}

		}

		if len(missingKeys) != 0 {
			return fmt.Errorf("Required keys missing or empty: %s", missingKeys)
		}
	}

	return nil
}

// RequiredKeys is a way for you to set a checkpoint when loading secrets or required exported variables for your application
func RequiredKeys(keys []string) {
	requiredKeys = append(requiredKeys, keys...)
}

// ApplyAdapter will set middleware, when Load or MustLoad is called those middleware will be called
func ApplyAdapter(a ...*Adapter) {
	adapters = append(adapters, a...)
}

// helper functions

func loadFiles(strict bool, filenames ...string) ([]string, error) {
	var files []string

	for _, filename := range filenames {
		f, err := os.Stat(filename)
		if err != nil {
			fmt.Printf("could not load %s: %s\n", filename, err)
			continue
		}

		if f.IsDir() {
			fmt.Printf("Could not load %s: %s", f.Name(), err)
			continue
		}

		bytes, err := ioutil.ReadFile(f.Name())
		if err != nil {
			return files, nil
		}

		files = append(files, string(bytes))
	}

	return files, nil
}

func setEnvMap(target *Map) error {
	for key, val := range target.Map {
		err := os.Setenv(key, val)
		if err != nil {
			return err
		}
	}

	return nil
}

// Parse takes a io.Reader that will parsed and returns a env map
func Parse(content string) *Map {
	emap := NewMap()

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		key, val := parseLine(line)

		if !strings.HasPrefix(key, "#") && key != "" {
			emap.Set(key, val)
		}
	}

	return emap
}

func parseLine(line string) (string, string) {
	trimed := strings.Trim(line, " ")

	splitLine := strings.Split(trimed, "=")

	return splitLine[0], splitLine[1]
}

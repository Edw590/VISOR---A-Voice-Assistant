/*******************************************************************************
 * Copyright 2023-2024 Edw590
 *
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 ******************************************************************************/

package Utils

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

/*
GPath (GoodPath) is sort of a copy of the string type but that represents a *surely* valid and correct path, also
according to the project conventions as described in the Path() function.

It's a "good path" because it's only given by Path(), which corrects the paths, and because the string component is
private to the package and only requested when absolutely necessary, like to communicate with Go's official functions
that require a string.
*/
type GPath struct {
	// p is the string that represents the path.
	p string
	// s is the path separator of the path
	s string
	// dir is true if the path *describes* a directory, false if it *describes* a file (means no matter if it exists and
	// we have permissions to read it or not).
	dir bool
}

type FileInfo struct {
	// Name is the name of the file
	Name       string
	// Modif_time is the last modification time of the file in Unix nanoseconds
	Modif_time int64
	// GPath is the path to the file
	GPath      GPath
}

/*
PathFILESDIRS combines a path from the given subpaths of type string or GPath (ONLY) into a GPath.

Note: the path separators used are always converted to the OS ones.

-----------------------------------------------------------

– Params:
  - separator – the path separator to use or "" for the OS default one
  - sub_paths – the subpaths to combine

– Returns:
  - the final path as a GPath
*/
func PathFILESDIRS(describes_dir bool, separator string, sub_paths ...any) GPath {
	var sub_paths_str []string = nil
	for _, sub_path := range sub_paths {
		val_str, ok := sub_path.(string)
		if ok {
			sub_paths_str = append(sub_paths_str, val_str)

			continue
		}

		val_GPath, ok := sub_path.(GPath)
		if ok {
			sub_paths_str = append(sub_paths_str, val_GPath.p)

			continue
		}

		// If it's not a string or GPath, it's an error.
		panic(errors.New("pathFILESDIRS() received an invalid type of parameter. " + getVariableInfoGENERAL(sub_path)))
	}

	if len(sub_paths_str) == 0 {
		return GPath{}
	}

	if describes_dir {
		// If the path describes a directory, make sure it ends with a path separator.
		if !strings.HasSuffix(sub_paths_str[len(sub_paths_str)-1], separator) {
			sub_paths_str[len(sub_paths_str)-1] += separator
		}
	}

	// Replace all the path separators with the OS path separator.
	for i, sub_path := range sub_paths_str {
		// Replace all the path separators with the OS path separator for Join() to work (only works with the OS one).
		sub_path = strings.Replace(sub_path, "/", string(os.PathSeparator), -1)
		sub_path = strings.Replace(sub_path, "\\", string(os.PathSeparator), -1)

		sub_paths_str[i] = sub_path
	}

	// Check if the last subpath ends in a path separator before calling Join() which will remove it if it's there.
	var ends_in_separator bool = false
	if strings.HasSuffix(sub_paths_str[len(sub_paths_str)-1], string(os.PathSeparator)) {
		ends_in_separator = true
	}

	if "" == separator {
		separator = string(os.PathSeparator)
	}

	// The call to Join() is on purpose - it correctly joins *and cleans* the final path string (only if it's used with
	// the OS path separator - which is always the case).
	var gPath GPath = GPath{
		p:   strings.Replace(filepath.Join(sub_paths_str...), string(os.PathSeparator), separator, -1),
		s:   separator,
		dir: false,
	}
	gPath.dir = gPath.DescribesDir()

	// Check if the path describes a directory and if it does, make sure the path separator is at the end (especially
	// since Join() removes it if it's there).
	if gPath.Exists() {
		if gPath.dir && !strings.HasSuffix(gPath.p, gPath.s) {
			gPath.p += gPath.s
		}
	} else {
		// As last resort, check through the last character on the subpaths list (project convention).
		if ends_in_separator && !strings.HasSuffix(gPath.p, gPath.s) {
			gPath.p += gPath.s
		}
	}

	return gPath
}

/*
Add adds subpaths to a path using the given path separator.

-----------------------------------------------------------

– Params:
  - separator – the path separator to use or 0xFF for the OS default one
  - sub_paths – the subpaths to add

– Returns:
  - the final path as a GPath
*/
func (gPath GPath) Add(describes_dir bool, separator rune, sub_paths ...any) GPath {
	// Create a temporary slice with the first element + the subpaths, all in a 1D slice and all as the 1st parameter of
	// the Path function.
	var tmp []any = append([]any{gPath}, sub_paths...)

	var separator_tmp string = string(separator)
	if separator == 0xFF {
		separator_tmp = ""
	}
	return PathFILESDIRS(describes_dir, separator_tmp, tmp...)
}

/*
Add2 is a wrapper of Add() using the OS default path separator.

-----------------------------------------------------------

– Params:
  - sub_paths – the subpaths to add

– Returns:
  - the final path as a GPath
 */
func (gPath GPath) Add2(describes_dir bool, sub_paths ...any) GPath {
	return gPath.Add(describes_dir, 0xFF, sub_paths...)
}

/*
GPathToStringConversion converts a GPath to a string.

The function has a big name on purpose to discourage its use - only use it when absolutely necessary (for the reason
written in the GPath type), like to call Go official file/directory functions.

-----------------------------------------------------------

– Params:
  - path – the path to convert

– Returns:
  - the GPath
*/
func (gPath GPath) GPathToStringConversion() string {
	return gPath.p
}

/*
ReadTextFile reads the contents of a file.

Note: all line breaks are replaced by "\n" for internal use, just like Python does.

-----------------------------------------------------------

– Returns:
  - the contents of the file or nil if an error occurs (including if the path describes a directory)
*/
func (gPath GPath) ReadTextFile() *string {
	if gPath.dir || !gPath.Exists() {
		return nil
	}

	data, err := os.ReadFile(gPath.p)
	if nil != err {
		return nil
	}
	var ret string = string(data)

	ret = strings.ReplaceAll(ret, "\r\n", "\n")
	ret = strings.ReplaceAll(ret, "\r", "\n")

	return &ret
}

/*
ReadFile reads the raw contents of a file.

-----------------------------------------------------------

– Returns:
  - the raw contents of the file or nil if an error occurs (including if the path describes a directory)
*/
func (gPath GPath) ReadFile() []byte {
	if gPath.dir || !gPath.Exists() {
		return nil
	}

	data, err := os.ReadFile(gPath.p)
	if nil != err {
		return nil
	}

	return data
}

/*
WriteTextFile writes the contents of a text file, creating it and any directories if necessary.

Note: all line breaks are replaced with the OS line break(s). So for Windows, "\r" and "\n" are replaced with "\r\n" and
for any other, "\r\n" and "\r" are replaced by "\n".

-----------------------------------------------------------

– Params:
  - content – the contents to write

– Returns:
  - nil if the file was written successfully, an error otherwise (including if the path describes a directory)
*/
func (gPath GPath) WriteTextFile(content string, append bool) error {
	var new_content string = content
	if "windows" == runtime.GOOS {
		new_content = strings.ReplaceAll(new_content, "\r\n", "\n")
		new_content = strings.ReplaceAll(new_content, "\r", "\n")
		new_content = strings.ReplaceAll(new_content, "\n", "\r\n")
	} else {
		new_content = strings.ReplaceAll(new_content, "\r\n", "\n")
		new_content = strings.ReplaceAll(new_content, "\r", "\n")
	}

	return gPath.WriteFile([]byte(new_content), append)
}

/*
WriteFile writes the raw contents of a file, creating it and any directories if necessary.

-----------------------------------------------------------

– Params:
  - content – the contents to write

– Returns:
  - nil if the file was written successfully, an error otherwise (including if the path describes a directory)
 */
func (gPath GPath) WriteFile(content []byte, append bool) error {
	if gPath.dir || gPath.Create(true) != nil {
		return errors.New("the path describes a directory or it couldn't be created")
	}

	if append {
		var flags int = os.O_WRONLY | os.O_APPEND

		file, err := os.OpenFile(gPath.p, flags, 0o777)
		if nil != err {
			return err
		}
		defer file.Close()

		_, err = file.Write(content)
	} else {
		// This way is here too because the other one doesn't seem to work well when it's not to append. Sometimes it
		// adds stuff to the file who knows why. This way here doesn't at least.
		err := os.WriteFile(gPath.p, content, 0o777)
		if nil != err {
			return err
		}
	}

	// Set the permissions to 777 after writing the file to be sure the file is accessible (OpenFile only sets the
	// permissions for the file creation).
	_ = os.Chmod(gPath.p, 0o777)

	return nil
}

/*
DescribesDir checks if a path DESCRIBES a directory or a file - means no matter if it exists or we have permissions to
see it or not.

It first checks if the path exists and if it does, checks if it's a directory or not - else it resorts to the path
string only, using the project convention in which a path that ends in a path separator is a directory and one that
doesn't is a file.

-----------------------------------------------------------

– Returns:
  - true if the path describes a directory, false if it describes a file
 */
func (gPath GPath) DescribesDir() bool {
	file_info, err := os.Stat(gPath.p)
	if err == nil {
		return file_info.IsDir()
	}

	return strings.HasSuffix(gPath.GPathToStringConversion(), gPath.s)
}

/*
Exists checks if a path exists.

-----------------------------------------------------------

– Returns:
  - true if the path exists (meaning the program also has permissions to *see* the file), false otherwise
*/
func (gPath GPath) Exists() bool {
	if nil != gPath.IsSupported() {
		return false
	}

	_, err := os.Stat(gPath.p)
	return err == nil
}

/*
Create creates a path and any necessary subdirectories in case they don't exist already.

-----------------------------------------------------------

– Params:
  - create_file – if true, creates the file *too* if the path represents a file

– Returns:
  - nil if the path was created successfully, an error otherwise
*/
func (gPath GPath) Create(create_file bool) error {
	if err := gPath.IsSupported(); nil != err {
		return err
	}

	var path_list []string = strings.Split(gPath.p, gPath.s)
	var describes_file bool = false
	if !gPath.dir {
		// If the path is a file, remove the file part of the file from the list so that it describes a directory only,
		// but memorize if it describes a file.
		describes_file = true
		path_list = path_list[:len(path_list) - 1]
	}

	// Create all parent directories if they don't exist.
	if !PathFILESDIRS(true, "", gPath.p[:FindAllIndexesGENERAL(gPath.p, gPath.s)[len(path_list)-1]+1]).Exists() {
		var current_path GPath = GPath{}
		if strings.HasPrefix(gPath.p, gPath.s) {
			current_path.p = gPath.s
		}
		for _, sub_path := range path_list {
			if "" == sub_path {
				continue
			}

			// Keep adding the subpaths until we reach the file part of the path, where the loop stops.
			current_path.p += sub_path + gPath.s

			if !current_path.Exists() {
				if err := os.Mkdir(current_path.p, 0o777); err == nil {
					_ = os.Chmod(current_path.p, 0o777)
				} else {
					return err
				}
			}
		}
	}

	// Create the file if the path represents a file.
	if create_file && describes_file && !gPath.Exists() {
		file, err := os.Create(gPath.p)
		if nil != err {
			return err
		}
		_ = os.Chmod(gPath.p, 0o777)
		_ = file.Close()
	}

	return nil
}

/*
Remove removes a file or directory.

-----------------------------------------------------------

– Returns:
  - nil if the file or directory was removed successfully, an error otherwise
 */
func (gPath GPath) Remove() error {
	if err := gPath.IsSupported(); nil != err {
		return err
	}

	return os.Remove(gPath.p)
}

/*
IsSupported checks if the path is supported by the current OS.

-----------------------------------------------------------

– Returns:
  - nil if the path is supported by the current OS, an error otherwise
*/
func (gPath GPath) IsSupported() error {
	// If the path is relative, it works everywhere (it's not specific to any OS). If it's absolute, it's supported
	// if it's an absolute path for the current OS.

	// Note: can't check with filepath.IsAbs() because it returns false for paths that are not supported by the current
	// OS, but are absolute for another OS. So the check must be made manually.

	// Check if the path is the wrong absolute type for the current OS. Else it's supported.
	// Don't forget the separators are changed to the current OS ones, so the checks are "inverted".
	if "windows" == runtime.GOOS {
		// Then check if it's a Linux absolute path.
		if strings.HasPrefix(gPath.p, "\\") {
			return errors.New("the path is not supported by the current OS")
		}
	} else {
		// Then check if it's a Windows absolute path.
		if len(gPath.p) >= 2 && (((gPath.p[0] >= 'a' && gPath.p[0] <= 'z' || gPath.p[0] >= 'A' && gPath.p[0] <= 'Z') && gPath.p[1] == ':') ||
					strings.HasPrefix(gPath.p, "//")) {
			return errors.New("the path is not supported by the current OS")
		}
	}
	// Else it's relative or absolute for the current OS.

	return nil
}

/*
GetFileList gets the list of files in a directory.

-----------------------------------------------------------

– Returns:
  - the list of files in the directory or nil if the path describes a file or an error occurs
 */
func (gPath GPath) GetFileList() []FileInfo {
	if !gPath.dir {
		return nil
	}

	files, err := os.ReadDir(gPath.GPathToStringConversion())
	if nil != err {
		return nil
	}

	var files_to_send []FileInfo = make([]FileInfo, 0, len(files))
	for _, file := range files {
		var file_path GPath = gPath.Add2(false, file.Name())
		file_stats, _ := os.Stat(file_path.GPathToStringConversion())
		files_to_send = append(files_to_send, FileInfo{
			Name:       file.Name(),
			Modif_time: file_stats.ModTime().UnixNano(),
			GPath:      file_path,
		})
	}

	return files_to_send
}

/*
GetOldestFileFILESDIRS gets the oldest file from a list of files.

-----------------------------------------------------------

– Params:
  - files_info – the list of files

– Returns:
  - the oldest file or an empty FileInfo if the list is empty
  - the index of the oldest file or -1 if the list is empty
 */
func GetOldestFileFILESDIRS(files_info []FileInfo) (FileInfo, int) {
	if len(files_info) == 0 {
		return FileInfo{}, -1
	}

	var oldest_file FileInfo = files_info[0]
	var oldest_file_idx int = 0
	for i := 1; i < len(files_info); i++ {
		if files_info[i].Modif_time < oldest_file.Modif_time {
			oldest_file = files_info[i]
			oldest_file_idx = i
		}
	}

	return oldest_file, oldest_file_idx
}

/*
GetBinDirFILESDIRS gets the full path to the directory of the binaries.

-----------------------------------------------------------

– Returns:
  - the full path to the directory of the binaries
*/
func GetBinDirFILESDIRS() GPath {
	return PersonalConsts_GL._VISOR_DIR.Add2(true, _BIN_REL_DIR)
}

/*
GetWebsiteFilesDirFILESDIRS gets the full path to the website files directory.

-----------------------------------------------------------

– Returns:
  - the full path to the website files directory
*/
func GetWebsiteFilesDirFILESDIRS() GPath {
	return PersonalConsts_GL._WEBSITE_DIR.Add2(true, _WEBSITE_FILES_REL_DIR)
}

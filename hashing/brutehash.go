package main

/*
 * [Example]:
 *     $ go build main.go 
 *     $ ./main md5 dictionary.txt 5d41402abc4b2a76b9719d911017c592
 *     $ ./main sha256 dictionary.txt 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
 *     $ ./main sha512 dictionary.txt 9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043
 */

import (
    "os"
    "fmt"
    "bufio"
    "strings"
    "encoding/hex"
)

import (
    "crypto/md5"
    "crypto/sha256"
    "crypto/sha512"
)

func main () {
    if len(os.Args) < 4 { get_error("args < 4") }
    check_mode(os.Args[1])
    hack_password(os.Args[1], os.Args[2], os.Args[3])
}

func check_mode (mode string) {
    if mode != "sha256" && mode != "sha512" && mode != "md5" {
        get_error("crypto-hash function is not found")
    }
}

func hack_password (mode string, dictionary, stat_hash string) {
    file, err := os.Open(dictionary)
    check_error(err)
    defer file.Close()

    var reader *bufio.Reader = bufio.NewReader(file)
    var pasw, dynamic_hash string

    fmt.Println("----------------------------------")
    for {
        pasw, _ = reader.ReadString('\n')
        pasw = strings.Replace(pasw, "\n", "", -1)
        dynamic_hash = encrypt(mode, pasw)
        if stat_hash == dynamic_hash {
            fmt.Println("----------------------------------")
            fmt.Println("[SUCCESS]:", pasw)
            fmt.Println("----------------------------------")
            os.Exit(0)
        } else {
            fmt.Println("[FAILURE]:", pasw)
        }
    }
}

func encrypt (crypt string, text string) string {
    if crypt == "md5" { 
        hash := md5.New() 
        hash.Write([]byte(text))
        return hex.EncodeToString(hash.Sum(nil))
    } else if crypt == "sha256" { 
        hash := sha256.New() 
        hash.Write([]byte(text))
        return hex.EncodeToString(hash.Sum(nil))
    } else { 
        hash := sha512.New() 
        hash.Write([]byte(text))
        return hex.EncodeToString(hash.Sum(nil))
    } 
}

func check_error (err error) {
    if err != nil { 
        fmt.Println("Error:", err)
        os.Exit(1) 
    }
}

func get_error (err string) {
    fmt.Println("Error:", err)
    os.Exit(1)
}

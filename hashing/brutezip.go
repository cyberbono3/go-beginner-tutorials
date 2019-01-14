package main

/*
 *  [Example]:
 *      $ go build main.go 
 *      $ ./main archive.zip dictionary.txt 
 */

import (
    "os"
    "fmt"
    "time"
    "bufio"
    "os/exec"
    "strings"
    "runtime"
)

func main () {
    var numcpu int = runtime.NumCPU()
    runtime.GOMAXPROCS(numcpu)

    check_args(os.Args)
    hack_password(numcpu, os.Args[1], os.Args[2])
}

func hack_password (numcpu int, archive, dictionary string) {
    file, err := os.Open(dictionary)
    check_error(err)
    defer file.Close()

    var reader *bufio.Reader = bufio.NewReader(file)
    fmt.Println("----------------------------------")

    var (
        index int
        pasw  string
    )

    for index = 0; index < numcpu; index++ {
        go func() {
            for {
                pasw, _ = reader.ReadString('\n')
                extract(archive, pasw[:len(pasw)-1])
            }
        } ()
        time.Sleep(time.Millisecond * 50)
    }

    fmt.Scanln()
}

func extract (archive, password string) {
    var commandString string = fmt.Sprintf("7z x %s -p%s -oExtractArchive -aoa", archive, password)
    var commandSlice []string = strings.Split(commandString, " ")

    var err error = exec.Command(commandSlice[0], commandSlice[1:]...).Run()
    check_result(err, password)
}

func get_error (err string) {
    fmt.Println("Error:", err)
    os.Exit(1)
}

func check_args (args []string) {
    if (len(args) < 3) {
        get_error("args < 3")
    }
    check_exist(args[1])
}

func check_exist (archive string) {
    var err error
    _, err = os.Stat(archive)
    check_error(err)
}

func check_error (err error) {
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
}

func check_result (err error, password string) {
    if (err == nil) {
        fmt.Println("----------------------------------")
        fmt.Printf("[SUCCESS]: %s\n", password)
        fmt.Println("----------------------------------")
        os.Exit(0)
    } else {
        fmt.Printf("[FAILURE]: %s\n", password)
    }
}

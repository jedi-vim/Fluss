package utils

import (
	"fmt"
	"time"
)

type ProgressStatus struct{
    Message string
    Finished bool
}

func AnimationBar(fn func(chan ProgressStatus)){
    finishedChannel := make(chan ProgressStatus)
    go fn(finishedChannel)

    status :=  <- finishedChannel
    for status.Finished == false{
        fmt.Printf("\r%s", status.Message)
        time.Sleep(500 * time.Millisecond)
        status = <- finishedChannel
    }
}

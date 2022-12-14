// region: packages

package rasa

import (
	"errors"
)

// endregion: packages
// region: messages

var (
	ErrInvalidSubCmd = errors.New("there must be a subcommand")
)

// endregion: messages

/*
func Exec(r Rasa, subCmd []string, data []byte) (result []byte, err error) {

	// subcommand

	Logger.Out(log.LOG_DEBUG, "rasaCmd Exec() subCmd", subCmd)

	if len(subCmd) == 0 {
		err = ErrInvalidSubCmd
		Logger.Out(log.LOG_ERR, err)
		return
	}
	if len(subCmd[0]) == 0 {
		err = ErrInvalidSubCmd
		Logger.Out(log.LOG_ERR, err)
		return
	}

	// data

	Logger.Out(log.LOG_DEBUG, "rasaCmd Exec() data", data)

	// endregion: validations
	// region: compile w/ ins and outs

	cmd := exec.Command(rasaCmd, subCmd...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		Logger.Out(log.LOG_ERR, err)
		return
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		Logger.Out(log.LOG_ERR, err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		Logger.Out(log.LOG_ERR, err)
		return
	}

	// endregion: compile
	// region: exec

	if err = instance.Lock(c); err != nil {
		Logger.Out(log.LOG_ERR, err)
		return
	}

	data = append(data, '\n')
	cmd.Start()
	stdin.Write(data)
	stdin.Close()

	// endregion: exec
	// region: scan stderr

	scanErr := bufio.NewScanner(stderr)
	for scanErr.Scan() {
		Logger.Out(log.LOG_WARNING, "rasa.Exec()", scanErr.Text())
	}
	if err := scanErr.Err(); err != nil {
		Logger.Out(log.LOG_ERR, err)
	}

	// endregion: stderr
	// region: scan stdout

	scanOut := bufio.NewScanner(stdout)
	for scanOut.Scan() {
		result = append(result, scanOut.Bytes()...)
		result = append(result, '\n')
	}
	Logger.Out(log.LOG_NOTICE, "rasaCmd.Exec() result", "\n"+string(result))
	if err = scanOut.Err(); err != nil {
		Logger.Out(log.LOG_ERR, err)
		instance.Unlock(c)
		return
	}
	if err = cmd.Wait(); err != nil {
		Logger.Out(log.LOG_ERR, err)
		instance.Unlock(c)
		return
	}

	// endregion: stdout

	instance.Unlock(c)
	return
}

*/

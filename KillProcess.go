package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	targetUser = "gleb"
)

func main() {
	fmt.Printf("Запуск программы для контроля сеансов пользователя %s...\n", targetUser)

	for {
		currentTime := time.Now()
		currentHour := currentTime.Hour()

		isForbiddenTime := false

		// Проверка на запрещенное время
		if (currentHour >= 21 && currentHour <= 23) || (currentHour >= 0 && currentHour < 18) {
			isForbiddenTime = true
		}

		if isForbiddenTime {
			fmt.Printf("Текущее время %s попадает в запрещенный интервал (с 21:00 до 18:00). Начинаю проверку сеансов пользователя %s.\n", currentTime.Format("15:04:05"), targetUser)
			terminateUserSessions(targetUser)
		} else {
			fmt.Printf("Текущее время %s находится в разрешенном интервале. Сеансы пользователя %s не будут завершены.\n", currentTime.Format("15:04:05"), targetUser)
		}

		time.Sleep(time.Minute)
	}

}

func terminateUserSessions(username string) {
	var cmd *exec.Cmd
	//var output []byte
	var err error

	switch runtime.GOOS {
	case "linux":
		// На Linux используем pkill для завершения процессов пользователя
		// -u <user> : процессы, принадлежащие указанному пользователю
		// -KILL     : отправить сигнал KILL (безусловное завершение)
		cmd = exec.Command("pkill", "-u", username, "-KILL")
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Ошибка при попытке завершить сеансы пользователя %s на Linux: %v\n", username, err)
			return
		}
		fmt.Printf("Успешно завершены сеансы пользователя %s на Linux (pkill -u %s -KILL).\n", username, username)

	case "darwin": // macOS
		// На macOS также можно использовать pkill, но иногда полезно явно перечислить и завершить процессы
		// Обойтись pkill проще и универсальнее.
		cmd = exec.Command("pkill", "-u", username, "-KILL")

		/* output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Ошибка при попытке завершить сеансы пользователя %s на macOS: %v\n", username, err)
			fmt.Printf("Вывод команды: %s\n", string(output))
			return
		} */

		err := cmd.Run()
		if err != nil {
			fmt.Printf("Ошибка при попытке завершить сеансы пользователя %s на macOS: %v\n", username, err)
			//fmt.Printf("Вывод команды: %s\n", string(output))
			return
		}

		fmt.Printf("Успешно завершены сеансы пользователя %s на macOS (pkill -u %s -KILL).\n", username, username)

	default:
		fmt.Printf("Операционная система %s не поддерживается.\n", runtime.GOOS)
		return
	}
}

// Вспомогательная функция для получения PID'ов процессов пользователя (не используется в текущей версии, но может быть полезна)
func getPIDsForUser(username string) ([]string, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux", "darwin":
		// ps -eo pid,user --no-headers | awk '$2 == "username" {print $1}'
		cmd = exec.Command("bash", "-c", fmt.Sprintf("ps -eo pid,user --no-headers | awk '$2 == \"%s\" {print $1}'", username))
	default:
		return nil, fmt.Errorf("операционная система %s не поддерживается для получения PID'ов", runtime.GOOS)
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения команды ps: %v", err)
	}

	pids := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(pids) == 1 && pids[0] == "" { // Если нет процессов
		return []string{}, nil
	}
	return pids, nil
}

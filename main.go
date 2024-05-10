package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Введите название файла")
	}
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer file.Close()
	text := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text = append(text, scanner.Text())

	}

	n, err := strconv.Atoi(text[0])
	if err != nil {
		fmt.Println("Error: ", err)
	}

	Tables := make([]Table, n)
	TimeSlice := strings.Split(text[1], " ")
	StartTime := ParseTime(TimeSlice[0]) // Время открытия клуба
	EndTime := ParseTime(TimeSlice[1])   // Время закрытия клуба

	PriceForHour, err := strconv.Atoi(text[2]) // Цена стола за час
	if err != nil {
		fmt.Println("Error: ", err)
	}
	ClientsAtTables := make(map[string]int, 0) //Для отслеживания клиентов за столом
	ClientsInClub := make(map[string]bool, 0)  // Для отслеживания клиентов в клубе
	queue := make([]string, 0)                 // Очередь клиентов
	fmt.Println(StartTime.Format("15:04"))
	for i := 3; i < len(text); i++ {
		fmt.Println(text[i])
		DataSlice := strings.Split(text[i], " ")
		action := DataSlice[1]
		switch action {
		case "1":
			ClientTime := ParseTime(DataSlice[0])
			if ClientTime.Before(StartTime) || ClientTime.After(EndTime) { // Проверяем, пришел ли клиент в часы работы клуба
				fmt.Println(DataSlice[0], "13 NotOpenYet")
				continue
			}
			if ClientsInClub[DataSlice[2]] { // Проверяем наличие клиента в клубе
				fmt.Println(DataSlice[0], "13 YouShallNoPass")
				continue
			} else {
				ClientsInClub[DataSlice[2]] = true
			}

		case "2":
			StartTime := ParseTime(DataSlice[0])
			if !ClientsInClub[DataSlice[2]] { //Проверяем наличие клиента в клубе
				fmt.Println(DataSlice[0], "13 ClientUnknown")
				continue
			}
			NumberOfTable, err := strconv.Atoi(DataSlice[3])
			NumberOfTable -= 1
			if err != nil {
				fmt.Println("Error: ", err)
			}
			if Tables[NumberOfTable].IsBusy { // Проверяем, свободен ли стол
				fmt.Println(DataSlice[0], "13 PlacesIsBusy")
				continue
			}
			if _, ok := ClientsAtTables[DataSlice[2]]; ok { // Случай, когда клиент захотел пересесть
				PreviousChair := ClientsAtTables[DataSlice[2]]
				EndTimeTable := ParseTime(DataSlice[0])
				Tables[PreviousChair].CountingCostAndTime(EndTimeTable, PriceForHour)
				Tables[PreviousChair].IsBusy = false
			}
			Tables[NumberOfTable].TakeTheTable(StartTime)
			ClientsAtTables[DataSlice[2]] = NumberOfTable
		case "3":
			flag := 0
			for i := range Tables {
				if !Tables[i].IsBusy { //Проверяем, есть ли свободный стол
					fmt.Println(DataSlice[0], "13 ICanWaitNoLonger!")
					flag = 1
					break
				}
			}
			if flag == 1 {
				continue
			}
			if len(queue) > n { // Случай, когда очередь больше чем количество столов
				fmt.Println(DataSlice[0], "11", DataSlice[2])
				continue
			} else {
				queue = append(queue, DataSlice[2])
			}
		case "4":
			if !ClientsInClub[DataSlice[2]] { // Проверяем наличие клиента в клубе
				fmt.Println(DataSlice[0], "13 ClientUnknown")
				continue
			}
			NumberOfTable := ClientsAtTables[DataSlice[2]]
			EndTimeTable := ParseTime(DataSlice[0])
			Tables[NumberOfTable].CountingCostAndTime(EndTimeTable, PriceForHour)
			Tables[NumberOfTable].IsBusy = false
			delete(ClientsAtTables, DataSlice[2])
			delete(ClientsInClub, DataSlice[2])
			if len(queue) > 0 { // Если в очереди кто то стоит, садим его за освободившийся стол
				firstClientInQueue := queue[0]
				fmt.Println(DataSlice[0], "12", firstClientInQueue, NumberOfTable+1)
				ClientsAtTables[firstClientInQueue] = NumberOfTable
				Tables[NumberOfTable].TakeTheTable(EndTimeTable)
				queue = queue[1:]
			}
		}
	}
	if len(ClientsAtTables) > 0 { // Узнаем оставшихся клиентов до закрытия клуба
		clientsAfterEndTime := make([]string, 0)
		for k, v := range ClientsAtTables {
			Tables[v].CountingCostAndTime(EndTime, PriceForHour)
			clientsAfterEndTime = append(clientsAfterEndTime, k)
		}
		sort.Strings(clientsAfterEndTime)
		for i := range clientsAfterEndTime {
			fmt.Println(EndTime.Format("15:04"), "11", clientsAfterEndTime[i])
		}
	}
	fmt.Println(EndTime.Format("15:04"))
	for i := range Tables { // Выводим выручку каждого стола за день работы, время проведенное за каждым столом
		fmt.Println(i+1, Tables[i].CostPerDay, Tables[i].TotalTime.Format("15:04"))
	}
}

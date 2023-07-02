package main

import "fmt"

func main() {
	// var ch = make(chan int): Khởi tạo một kênh có tên ch kiểu dữ liệu là int.
	// Kênh này được sử dụng để truyền dữ liệu từ goroutine đến main goroutine (goroutine chính).
	var ch = make(chan int)
	// close(ch): Đóng kênh ch. Điều này cho biết không còn dữ liệu nào sẽ được gửi vào kênh này
	// và bất kỳ lần gửi dữ liệu nào sau đó đều gây ra lỗi runtime.
	close(ch)

	var ch2 = make(chan int)
	// Đây là một goroutine mới được tạo ra.
	// Goroutine này chạy một vòng lặp từ 1 đến 9 và gửi các giá trị vào kênh ch2.
	go func() {
		for i := 1; i < 10; i++ {
			ch2 <- i
		}
		close(ch2)
	}()

	// Dùng select để chọn nhận dữ liệu từ hai kênh khác nhau: ch và ch2.
	// Nếu nhận dữ liệu thành công từ ch, in ra "ch1", giá trị nhận được (x) và trạng thái kênh (ok).
	// Nếu ok là false, điều này có nghĩa là kênh ch đã bị đóng, và ta gán ch = nil để không nhận dữ liệu từ kênh này nữa.
	// Cuối cùng, kiểm tra nếu cả ch và ch2 đều bị đóng (nil), thoát khỏi vòng lặp và kết thúc chương trình.
	for { // Vong lap vo han
		select { // chọn 1 trong nhiều hoạt động (action) trên các kênh giao tiếp,
		// select là một câu lệnh blocking, nghĩa là nó sẽ dừng lại và chờ đợi cho tới khi một trong các
		// trường hợp trong danh sách các trường hợp của nó có thể thực hiện được
		// select sử dụng với câu lệnh `case` để liệt kê các trường hợp vmà chương trình muốn xử lý
		// Khi có một hoạt động có thể thực hiện được trên một trong các kênh trong danh sách trường hợp, hoạt động đó sẽ được thực hiện và khối lệnh tương ứng với case đó sẽ được thực thi.
		// action sẽ được thực hiện và khối lệnh tương ứng với `case` đó để thực thi
		// nếu có nhiều hơn một hoạt động có thể thực hiện được cùng một lúc, golang sẽ chọn ngẫu nhiên một trong số chúng để thực thi
		// nếu không có hoạt động nào được thực hiện và không có trường hợp `default` trong select, chương trình sẽ
		// bị block cho tới khi một hoạt động có thể thực hiện được
		// *Tóm lại, cấu trúc for được sử dụng để thực hiện vòng lặp,
		//trong khi cấu trúc select được sử dụng để chọn một trong nhiều hoạt động có thể thực hiện được từ các kênh giao tiếp.
		case x, ok := <-ch:
			fmt.Println("ch1", x, ok)
			if !ok {
				ch = nil
			}
		case x, ok := <-ch2:
			fmt.Println("ch2", x, ok)
			if !ok {
				ch2 = nil
			}
		}

		if ch == nil && ch2 == nil {
			break //thoát khỏi vòng lặp và kết thúc chương trình.
		}
	}
}

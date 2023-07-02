package internal

import (
	"errors"
	"fmt"
	"github.com/cheggaaa/pb"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func HandleHTTPRequests(reqs, results chan string, quit chan int, bar *pb.ProgressBar, details *RequestDetails) {

	for link := range reqs {

		log.Debug().Msg(link)

		client := http.Client{
			Transport: &http.Transport{
				DisableKeepAlives: true},
		}

		req, err := http.NewRequest("HEAD", "https://"+link, nil)

		if err != nil {
			results <- "err"
			bar.Increment()
			continue
		}

		if len(details.RandomAgent) > 0 {
			chosenAgent := SelectRandomItem(details.RandomAgent)
			req.Header.Set("User-Agent", chosenAgent)
		}

		resp, err := client.Do(req)

		if err != nil {

			results <- "err"
			bar.Increment()
			continue
		}

		//log.Debug().Msg(strconv.Itoa(resp.StatusCode))

		bar.Increment()
		results <- link + ":" + strconv.Itoa(resp.StatusCode)
	}

	if len(reqs) == len(results) {
		quit <- 0
	}

}

func AsyncHTTPHead(urls []string, threads int, timeout int, details RequestDetails, output string) {
	//Kênh result được khởi tạo để nhận kết quả từ các yêu cầu HTTP được xử lý.
	//Kênh reqs được khởi tạo với độ dài bằng với số lượng URL để chứa các URL cần gửi yêu cầu HTTP. Đây là một kênh có buffer.
	//Kênh quit được khởi tạo để thông báo khi tất cả các yêu cầu HTTP đã được xử lý xong.
	//Biến bar là một thanh tiến trình được tạo ra bằng thư viện pb để hiển thị tiến trình xử lý yêu cầu HTTP.
	result := make(chan string)
	reqs := make(chan string, len(urls)) // buffered
	quit := make(chan int)

	bar := pb.StartNew(len(urls))

	//Một vòng lặp for được sử dụng để khởi tạo các goroutine (luồng riêng biệt) để xử lý yêu cầu HTTP.
	//Số lượng goroutine được xác định bởi giá trị threads.
	for i := 0; i < threads; i++ {
		//Trong mỗi goroutine, hàm HandleHTTPRequests() được gọi để xử lý yêu cầu HTTP.
		//Các tham số như kênh reqs, kênh result, kênh quit, thanh tiến trình bar và con trỏ &details được truyền vào.
		go HandleHTTPRequests(reqs, result, quit, bar, &details)
	}

	//Một goroutine khác được tạo ra bằng cách sử dụng từ khóa go func() {...}()
	//để gửi các URL từ mảng urls vào kênh reqs.
	go func() {
		for _, link := range urls {
			fmt.Println(link)
			reqs <- link
		}
	}()

	//var results []string

	// parsing http codes
	// 500 , 502 server error
	// 404 not found
	// 200 found
	// 400, 401 , 403  protected
	// 302 , 301 redirect

	//Sau đó, chương trình thực hiện một vòng lặp vô hạn sử dụng câu lệnh for {}
	//và sử dụng câu lệnh select để xử lý các sự kiện từ các kênh:
	for {
		select {
		//Nếu có kết quả từ kênh result, nó được nhận và kiểm tra giá trị. Nếu giá trị khác "err",
		//chương trình phân tích và xử lý kết quả để tạo ra đầu ra tương ứng.
		//Đầu ra được ghi vào log và được ghi vào output.
		case res := <-result:
			fmt.Println(res)
			if res != "err" {
				domain := res
				var out, status string

				if strings.Contains(res, ":") {
					domain = strings.Split(res, ":")[0]
					status = strings.Split(res, ":")[1]
				}

				if status == "200" {
					out = fmt.Sprintf("%s: %s - %s", status, "Open", domain)
					log.Info().Msg(out)
				} else {
					out = fmt.Sprintf("%s: %s - %s", status, "Redirect", domain)
					log.Warn().Msg(out)
				}

				if out != "" {
					_, _ = AppendTo(output, out)
				}

			}

		//Nếu sau một khoảng thời gian chờ (timeout), không có kết quả nào từ kênh result,
		//chương trình ghi log để thông báo là "TimeOut" và tăng giá trị của thanh tiến trình bar.
		case <-time.After(time.Duration(timeout) * time.Second):
			log.Warn().Msg("TimeOut")
			bar.Increment()

		//Nếu nhận được thông báo từ kênh quit, chương trình cập nhật giá trị thanh tiến trình bar và kết thúc vòng lặp.
		case <-quit:
			bar.Set(len(urls))
			bar.Finish()
			return
		}
	}

}

func GenerateMutatedUrls(wordListPath string, mode string, provider string, providerPath string, target string, environments []string) ([]string, error) {

	//envs := []string{"test", "dev", "prod", "stage"}
	words, err := ReadTextFile(wordListPath)

	if err != nil {
		log.Fatal().Err(err).Msg("Exiting ...")
	}
	permutations := []string{"%s-%s-%s", "%s-%s.%s", "%s-%s%s", "%s.%s-%s", "%s.%s.%s"}

	var compiled []string

	for _, env := range environments {

		for _, word := range words {

			for _, permutation := range permutations {
				formatted := fmt.Sprintf(permutation, target, word, env)
				compiled = append(compiled, formatted)
			}

		}
	}

	urlPermutations := []string{"%s.%s", "%s-%s", "%s%s"}
	for _, word := range words {

		for _, permutation := range urlPermutations {
			formatted := fmt.Sprintf(permutation, target, word)
			compiled = append(compiled, formatted)
		}

	}

	providerConfig, err := InitCloudConfig(provider, providerPath)

	if err != nil {
		log.Fatal().Err(err).Msg("Exiting...")
	}

	log.Info().Msg("Initialized " + provider + " config")

	var finalUrls []string

	if mode == "storage" {

		if len(providerConfig.StorageUrls) < 1 && len(providerConfig.StorageRegionUrls) < 1 {
			return nil, errors.New("storage are not supported on :" + provider)
		}

		if len(providerConfig.StorageUrls) > 0 {

			for _, app := range providerConfig.StorageUrls {

				for _, word := range compiled {
					finalUrls = append(finalUrls, word+"."+app)
				}
			}
		}

		if len(providerConfig.StorageRegionUrls) > 0 {

			for _, region := range providerConfig.Regions {
				for _, regionUrl := range providerConfig.StorageRegionUrls {
					for _, word := range compiled {
						finalUrls = append(finalUrls, word+"."+region+"."+regionUrl)
					}
				}
			}
		}
	}

	if mode == "app" {
		if len(providerConfig.APPUrls) < 1 && len(providerConfig.AppRegionUrls) < 1 {
			return nil, errors.New("storage are not supported on :" + provider)
		}

		if len(providerConfig.APPUrls) > 0 {
			for _, app := range providerConfig.APPUrls {
				for _, word := range compiled {
					finalUrls = append(finalUrls, word+"."+app)
				}
			}
		}

		if len(providerConfig.AppRegionUrls) > 0 {
			for _, region := range providerConfig.Regions {
				for _, regionUrl := range providerConfig.AppRegionUrls {
					for _, word := range compiled {
						finalUrls = append(finalUrls, word+"."+region+"."+regionUrl)
					}
				}
			}
		}
	}

	return finalUrls, nil

}

## About
Serbian identification card reader.

Program to read serbian identification cards, aimed mainly at linux os.
You can read data from the cards and print it to pdf file.


## Instructions

Run it with: go run main.go genpdf.go.

Build: go build main.go genpdf.go 

## Notes

As stated [HERE][1] you have to install the PCSC-Lite daemon and CCID driver.

Cards issued before 2015 are not supported. 

If for some reason you **really** need it feel free to contact me.

[1]: https://github.com/sf1/go-card/


## Screen shots

![goSteel_sample_02](https://user-images.githubusercontent.com/37847620/212207417-a4b9aced-cc2e-487e-91fc-d8646c4c10e7.png)

![gosteel_sample_01](https://user-images.githubusercontent.com/37847620/211685384-4f6c1f00-5576-45bd-9934-c556c5491e2a.png)

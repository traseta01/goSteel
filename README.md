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

If for some reason you **really ** need it feel free to contact me.

[1]: https://github.com/sf1/go-card/

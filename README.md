## Parser Web Crawler

### Description

### Features
- crawl the page and follow the links of the same host
- robot.txt is taking in account when crawl (if exist)
- retry mechanism for too many request
- visit pages concurrently

### Project Structure

```
|_ /cmd (contains the main file) 
|_ /integration (contains the integration tests)
|_ /internal
    |_ /analyzer (logic for analyze pages)
    |_ /crawler  (core execution for crawling)
    |_ /fetcher  (retrieves the page info)
```

### Technologies used
- Golang 1.19
- Wiremock
- Testify library 
- Containers (Docker)


### Run the project
Makefile is provided to allow a management of the project.

### Run

```
./webcrawler -url={url} -depth={depth} -delay={delay}`

 -url: the url to crawl
 -depth: depth of the crawl
 -delay: retrieve pages with certain delay  



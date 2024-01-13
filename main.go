package main

import (
    "encoding/json"
    "github.com/gin-gonic/gin"
    "io/ioutil"
    "net/http"
    "fmt"
)

type GitHubUser struct {
    Login     string `json:"login"`
    AvatarURL string `json:"avatar_url"`
    Name      string `json:"name"`
    PublicRepos int   `json:"public_repos"`
    StarCount int `json:"star_count"`
}

func main() {
    r := gin.Default()
    r.LoadHTMLGlob("templates/*")

    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", nil)
    })

    r.POST("/fetch", func(c *gin.Context) {
        username := c.PostForm("username")
        user, err := fetchGitHubUser(username)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, user)
    })

    r.Run(":8080")
}

func fetchGitHubUser(username string) (*GitHubUser, error) {
    url := "https://api.github.com/users/" + username
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var user GitHubUser
    err = json.Unmarshal(body, &user)
    if err != nil {
        return nil, err
    }
    starCount, err := getUserStarCount(username)
    if err != nil {
        return nil, err
    }
    user.StarCount = starCount
    fmt.Println(starCount)

    return &user, nil
}



func getUserStarCount(username string) (int, error) {
    url := "https://api.github.com/users/" + username + "/repos"
    resp, err := http.Get(url)
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()

    var repos []struct {
        StargazersCount int `json:"stargazers_count"`
    }

    err = json.NewDecoder(resp.Body).Decode(&repos)
    if err != nil {
        return 0, err
    }

    starCount := 0
    for _, repo := range repos {
        starCount += repo.StargazersCount
    }
    fmt.Println(starCount)

    return starCount, nil
}

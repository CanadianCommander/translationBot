# Translation Bot 🤖
Slack bot to automate the translation process. 

### Usage 
TranslationBot is currently installed in our internal slack workspace. You can interact with TranslationBot 
in the following ways. 

#### Issuing slash commands 
- `/translation` - brings up help / command index. Any slash command can also be run by clicking the "Run" button
to the right of the command on this page. 
- `/translation projects` - lists all projects that you can use TranslationBot on
- `/translation missing <project>` - outputs missing translations by language for the specified project.

#### Uploading translation files 
TranslationBot will automatically detect translation files uploaded to the translations channel. TranslationBot
not only checks the file extension but also checks the format of the file. Currently, these file formats are supported.

*CSV_Tara_1.0*

This format is a CSV format. It is a CSV file with column headers, where each header is a language.
Each following row maps a string in each of the given languages. TranslationBot uses these mappings to 
apply translations by simply looking up the source string (ENGLISH). Then updates the other languages using 
the values provided in that row. See the following example.
```csv
ENGLISH, FRENCH, SPANISH, ANOTHER_LANGUAGE... 
Hello World, Bonjour le monde, Hola Mundo, .....
translation magic, traduction magique, traducción mágica, .....
.....
 ```

This format is not appropriate when one source string gos to multiple different strings in other languages.
To handle this you will have to use the following format

*CSV_Keyed_Tara_1.0*

This format is just like the CSV_Tara_1.0 format except, that the first column is `KEY`. This `KEY` is used to
disambiguate multiple source strings that map to the same string in another language. See the following example.
```csv
KEY, ENGLISH, FRENCH, SPANISH, ANOTHER_LANGUAGE...
login.title, Hello World, Bonjour le monde, Hola Mundo, .....
login.title.about, translation magic, traduction magique, traducción mágica, .....
.....
```
This is the preferred way of uploading translation files, as it reduces the chance of the wrong string getting updated. 


### Getting Started (Dev)

#### prerequisites 
- go 1.18.8+
- git  
- node 16+
- yarn 

#### setup 

Create a configuration file for TranslationBot under `./config/config.yaml` The file should be similar to 
`kubernetes/helm/templates/translation-bot-secrets.yaml`. It is up to you, which projects you configure. Just be 
sure to fill in the secret values such as `slackClientSecret` and `gitPassword` 

Now that you have a config file, simply start TranslationBot `go run ./cmd/tb` It's just that simple 

### Deployment 
TranslationBot is deployed via GitHub Actions CI/CD pipeline. Simply merge in to `main` and your changes 
automatically deploy to `api.cluster.bbenetti.ca`

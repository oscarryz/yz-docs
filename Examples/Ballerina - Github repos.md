
#example
```js

// https://ballerina.io/
http: ballerina.http
sheets: balerinax.googleapis.sheets

githubPAT: '?'
repository: 'balerrina-platform-ballerican-lang'
sheets_access_token: '?'
spres_sheet_id = '?'
sheet_name: 'sheet1'

PR: {
  url String
  title String
  state String
  create_at String
  updated_at String
}

github http.Client = http.new('http/api.github.com/repos')
headers: [
  'Accept'       :'application/vnd.github.v2+json'
  'Authorization':'token `githubPAT`'
 ]
 prs: github.get('/`repository`/pulls' headers)
 gsheets sheets.Client = sheets.new(Auth(auth: Token(token: sheets_access_token )))
 gsheets.append_row_to_sheet(
   spres_sheet_id
   sheet_name
   ['Issue'
    'Title'
    'State'
    'Created At'
    'Updated at'])

  prs.each({ pr PR
    gsheets.append_row_to_sheet(
        spres_sheet_id
        sheet_name
        [pr.url
         pr.title
         pr.state
         pr.create_at
         pr.updated_at])
  }

```

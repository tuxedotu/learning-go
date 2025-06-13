# litte bit of the web from back to front

## places
- [session-token & auth](https://www.sohamkamani.com/golang/session-cookie-authentication/)
- [golang html/template](https://pkg.go.dev/text/template#hdr-Functions)

## explore guide
- [x] favicon.ico
- [x] login/signup logic
- [x] htmx + go templating ('hx-'handles, {{}} mechanics, etc.)
- [x] action-based web flow
- [>] wrap up [session-token & auth](https://www.sohamkamani.com/golang/session-cookie-authentication/) 
    - [x] action: signup -> create user & sessionToken
    - [x] action: login -> update sessionToken
    - [ ] logic: remove token from cache after expiry (when?)
- [x] shared message board
    - collection of messages (author, msg, createdAt)
    - range in template
- [ ] *messageboard* -> live/automatic-client-side reload via htmx-trigger="load every Ns" or sth like it
    - make view hot-reloadable instead of (lazy ass) reroute
- [ ] *login/signup* -> add visual error ui
    - TODO: signup -> check for existing users
- [ ] **extra**: '/stats' route to see/print all current srv vars/caches/stores (client/srv-side?; over-the-air html updates?)

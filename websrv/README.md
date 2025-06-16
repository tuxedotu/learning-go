# litte bit of the web from back to front

## places
- [session-token & auth](https://www.sohamkamani.com/golang/session-cookie-authentication/)
- [golang html/template](https://pkg.go.dev/text/template#hdr-Functions)
- [font source (pixelgame)](https://www.fontspace.com/pixel-game-font-f121080)
- [gruvbox-css](https://github.com/YV31/gruvbox-css/tree/master)

## current todos:
- [X] *redirects* -> make 'default' into "/"-case and redirect to "/" by default 
- [ ] *ui*: -> **mobile** (input fields & login) || **navbar** (position)

## ideas to explore
- [x] favicon.ico
- [x] login/signup logic
- [x] htmx + go templating ('hx-'handles, {{}} mechanics, etc.)
- [x] action-based web flow
- [X] wrap up [session-token & auth](https://www.sohamkamani.com/golang/session-cookie-authentication/) 
    - [x] action: signup -> create user & sessionToken
    - [x] action: login -> update sessionToken
- [x] shared message board
    - collection of messages (author, msg, createdAt)
    - range in template
- [x] *messageboard* -> live/automatic-client-side reload via htmx-trigger="load every Ns" or sth like it
    - make view hot-reloadable instead of (lazy ass) reroute
- [x] *error ui* -> add visual error displays (login, sendMsg)
- [x] *signup* -> check for existing users
- [x] *messages* -> fix site-breaking overflow
- [ ] *keyboard interactivity* -> basic js-events & <ret> <tab> key inputs
- [ ] *custom icons* -> cursor, ux
    - fix submit-msg button icon
- [ ] *add color mechanic* -> press '1,2,3,...' for selection of username colors OR just random ??
    - add fields `User.Color` & `Message.Color` (for accurate history)
    - create `color-picker` template
    - add `color-picker` template to `login` site
- [ ] **extra**: '/stats' route to see/print all current srv vars/caches/stores (client/srv-side?; over-the-air html updates?)

## golang-notes
- **arrays / slices / maps** -> `[n]any` vs. `[]any` vs. `map[any]any`
    - *arrays*: fixed type-safe collection of data -> `arr := [2]int{1, 2}`
        - on empty init -> elements are set to type-specific zero-values
    - *slice*: non-fixed type-safe collection of data -> `arr := []int{1, 2}`
        - can be resized via builtin func `append(slice, ...any)` or `slice[len+1] = ...` 
        - on empty init -> elements are set to type-specific zero-values
    - *maps*: type-safe hash-map key-value pairs of data -> `map := map[int]string{1: "hello"}`
- **`make()` vs. `new()`** - modes of arr/map init
    - ??

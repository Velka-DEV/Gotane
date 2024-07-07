# Gotane: High performance credential stuffing library ![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/Velka-DEV/Gotane/go.yml?branch=main) [![Latest release](https://img.shields.io/github/v/release/Velka-DEV/Gotane?color=blue&style=flat-square&include_prereleases)](https://github.com/Velka-DEV/Gotane/releases) ![GitHub](https://img.shields.io/github/license/Velka-DEV/Gotane)

Gotane is a Go library for pentesting web apps against credential stuffing attacks library written by Velka.

In fact, it will manage:
* [x] Check process parallelization using the ants library (High performances)
* [x] Http clients creation and distribution (support for backconnect proxies and all protocols). The library make reuse of the clients for better performances (Thanks [@Laiteux](https://github.com/Laiteux))
* [x] Output (directory, whether to output invalids or not, console colors, captures and more...)
* [ ] Console title (hits, free hits, estimated hits, percentages, checked per minute, elapsed/remaining time and more...) (To be added)
* [x] Check pause/resume/end functions using hotkeys (Hotkeys to be added, actually only control methods are avaiable. Feel free to create yourself your bindings or make contribution to the project)

And more... See the code itself or the folder for a better overview.

## Credits

Thanks [@Laiteux](https://github.com/Laiteux) for the original inspiration with the [Milky](https://github.com/Laiteux/Milky) library. 

## Contribute

Your help and ideas are welcome, feel free to fork this repo and submit a pull request.

However, please make sure to follow the current code base/style.

## Contact

Telegram: @modest

## Donate

If you would like to support this project, please consider donating.

Donations are greatly appreciated and a motivation to keep improving.

- Bitcoin (BTC) 
    - Segwit: `bc1qhyzca03e78td68e2ppkqpdp6224pesw66vn6pv`
    - Legacy: `15xh18xUPL76wBYZ7qWCEwETQx6iPx8xTq`
- Monero (XMR): `44f7XFPeddmGDPvNDtR9sKQen619oVEGXaenw3XjKecWKrd4ZFtS6Md9yzrcYc3p47JVRnkQMFLvc8Eh8ua7n7D4BwNNjuY` 
- Litecoin (LTC): `LSnspYrX217dwRq7c9kd8hxuUkRL4K9z1N`
- Ethereum (ETH): `0xef20130259D5F3F094049a1eFb25f5c23052fDd8`

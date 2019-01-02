const domains = [
    "www.google.com"
];

/**
 * @return {string}
 */
function FindProxyForURL(url, host) {
    return "SOCKS5 127.0.0.1:2680";

    if (isInNet(host, "127.0.0.1", "255.0.0.0")) {
        return "DIRECT";
    }

    for (let i = 0; i < domains.length; i++) {
        if (domains[i] === url) {
            return "SOCKS5 127.0.0.1:2680; DIRECT;";
        }
    }

    return "DIRECT";
}
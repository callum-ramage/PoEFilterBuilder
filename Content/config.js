require.config({
    baseUrl: "https://web.poecdn.com/js/",
    paths : {"main":"main.3e3ca416463187bd8c66006318d480cd0f4926d9","plugins":"plugins.4bf418572b38a21b9b40f7510616448a71ab4f2c","skilltree":"skilltree.a971a0e90d0b8e73e7f33d78d232a97b9aa099f0","trade":"trade.82daf6321e091646ee3e38fc9164b8fa7e7008c3","itemfilter":"itemfilter.3d62cb8976ee009776b89e704ce7330297aaf010","notice":"notice.d4955e72932137ac9f94028dd14467ac6ee17c1f"},
    shim: {"main":{"deps":["config","plugins"]},"plugins":{"deps":["config"]},"notice":{"deps":["config"]}}
});

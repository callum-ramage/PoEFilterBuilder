require.config({
    baseUrl: "https://web.poecdn.com/js/",
    paths: { "main": "main.73a98a716924ec65b543d19174955642eab54412", "plugins": "plugins.fba22274dc3520c88a10bc27431e90615239367a", "skilltree": "skilltree.88a405683eaef0e980fad1ba8054a8e47e4ba89e", "trade": "trade.d9b5b0668845556379319a997f2b8b0432140338", "itemfilter": "itemfilter.d83b5cdad7ba76f6819bae03033e15e2a7212af5", "notice": "notice.d4955e72932137ac9f94028dd14467ac6ee17c1f" },
    shim: { "main": { "deps": ["config", "plugins"] }, "plugins": { "deps": ["config"] }, "notice": { "deps": ["config"] } }
});
document.addEventListener('DOMContentLoaded', function(event) {
    var navBar = document.getElementById("nav");
    createButton("Home", "/home", navBar);
    createButton("Filter", "/filter", navBar);
    createButton("Replace", "/replace", navBar);
    createButton("Configuration", "/config", navBar);
    createButton("Close", "/close", navBar);
});

function createButton(label, link, parent) {
    var button = document.createElement('BUTTON');
    button.onclick=function() {
        location.href=link;
    }
    var text = document.createTextNode(label); 
    button.appendChild(text); 
    parent.appendChild(button);
}
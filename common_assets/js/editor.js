function add_stylesheet(href){
  var link = document.createElement("LINK");
  link.setAttribute("rel", "stylesheet");
  link.setAttribute("href", href);
  document.head.appendChild(link)
}

function add_script(src){
  var script = document.createElement("SCRIPT");
  script.src = src;
  script.async = false
  script.defer = false

  document.head.appendChild(script)
}
add_stylesheet("common_assets/css/font-awesome.min.css")
add_script("common_assets/js/showdown.js")
add_script("common_assets/js/jquery-1.11.3.min.js")
add_script("common_assets/js/jquery-ui.js")
add_script("common_assets/js/editor-core.js")
add_stylesheet("common_assets/css/custom-theme/jquery-ui-1.10.3.custom.css")

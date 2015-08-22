var converter = new showdown.Converter();
function update_title(){
  $("#title").html($("#input_title").val());
}

function update_short_description(){
  $("#short_description").html($("#input_short_description").val());
}

function update(){
  var srcText = $('#src').val();
  var html = converter.makeHtml(srcText);
  $("#target").html(html);
}

var lastSelectionEnd = 0
$(document).ready(function() {
  $("#btnMD_H1").on('click', function() {
    var caretPos = document.getElementById("src").selectionEnd;
    var textAreaTxt = jQuery("#src").val();

    var i = caretPos;
    while (i >= 0) {
      i = i - 1
      if (textAreaTxt.charAt(i)=="\n") {
        i = i + 1
        break
      }
    }
    var firstPart = textAreaTxt.substring(0, i)
    var secondPart = textAreaTxt.substring(i)

    $("#src").val(firstPart + "#"+secondPart.replace(/^[\#>]*/, ''));
    update();
  });

  $("#btnMD_H2").on('click', function() {
    var caretPos = document.getElementById("src").selectionEnd;
    var textAreaTxt = jQuery("#src").val();

    var i = caretPos;
    while (i >= 0) {
      i = i - 1
      if (textAreaTxt.charAt(i)=="\n") {
        i = i + 1
        break
      }
    }
    var firstPart = textAreaTxt.substring(0, i)
    var secondPart = textAreaTxt.substring(i)

    $("#src").val(firstPart + "##"+secondPart.replace(/^[\#>]*/, ''));
    update();
  });

  $("#btnMD_H3").on('click', function() {
    var caretPos = document.getElementById("src").selectionEnd;
    var textAreaTxt = jQuery("#src").val();

    var i = caretPos;
    while (i >= 0) {
      i = i - 1
      if (textAreaTxt.charAt(i)=="\n") {
        i = i + 1
        break
      }
    }
    var firstPart = textAreaTxt.substring(0, i)
    var secondPart = textAreaTxt.substring(i)

    $("#src").val(firstPart + "###"+secondPart.replace(/^[\#>]*/, ''));
    update();
  });

  $("#btnMD_quote").on('click', function() {
    var caretPos = document.getElementById("src").selectionEnd;
    var textAreaTxt = jQuery("#src").val();

    var i = caretPos;
    while (i >= 0) {
      i = i - 1
      if (textAreaTxt.charAt(i)=="\n") {
        i = i + 1
        break
      }
    }
    var firstPart = textAreaTxt.substring(0, i)
    var secondPart = textAreaTxt.substring(i)

    $("#src").val(firstPart + ">"+secondPart.replace(/^[\#>]*/, ''));
    update();
  });

  $("#btnMD_list").on('click', function() {
    var textAreaTxt = jQuery("#src").val();

    var ini = document.getElementById("src").selectionStart
    var end = document.getElementById("src").selectionEnd
    var lines = jQuery("#src").val().split('\n');
    var newLines = new Array();

    var lineIni = 0
    for (i = 0; i < lines.length; i++) {
      var f = ""
      var l = ""
      lineEnd = lineIni + lines[i].length + 1

      var add = false
      if (ini >= lineIni && ini < lineEnd){
        if (i > 0 && lines[i-1].length > 1){
          f = "\n"
        }
        add = true
      }

      if (lineIni >= ini  && lineEnd <= end ){
        add = true
      }

      if (end >= lineIni && end < lineEnd){
        if (i + 1 < lines.length && lines[i+1].length > 1){
          l = "\n"
        }
        add = true
      }

      if (add) {
        newLines[i] = f+"- "+lines[i]+l
      } else {
        newLines[i] = lines[i]
      }

      lineIni = lineIni + lines[i].length + 1
    }

    var output = ""
    for (i = 0; i < newLines.length; i++) {
      output = output + newLines[i]
      if (i < newLines.length-1){
        output = output +"\n"
      }
    }
    $("#src").val(output);

    update();
  });

  $("#btnMD_orderedlist").on('click', function() {
    var textAreaTxt = jQuery("#src").val();

    var ini = document.getElementById("src").selectionStart
    var end = document.getElementById("src").selectionEnd
    var lines = jQuery("#src").val().split('\n');
    var newLines = new Array();

    var lineIni = 0
    var lineNo = 1
    for (i = 0; i < lines.length; i++) {
      var f = ""
      var l = ""
      lineEnd = lineIni + lines[i].length + 1

      var add = false
      if (ini >= lineIni && ini < lineEnd){
        if (i > 0 && lines[i-1].length > 1){
          f = "\n"
        }
        add = true
      }

      if (lineIni >= ini  && lineEnd <= end ){
        add = true
      }

      if (end >= lineIni && end < lineEnd){
        if (i + 1 < lines.length && lines[i+1].length > 1){
          l = "\n"
        }
        add = true
      }

      if (add) {
        newLines[i] = f+lineNo+". "+lines[i]+l
        lineNo = lineNo+1
      } else {
        newLines[i] = lines[i]
      }

      lineIni = lineIni + lines[i].length + 1
    }

    var output = ""
    for (i = 0; i < newLines.length; i++) {
      output = output + newLines[i]
      if (i < newLines.length-1){
        output = output +"\n"
      }
    }
    $("#src").val(output);

    update();
  });

  $("#btnMD_code").on('click', function() {
    var textAreaTxt = jQuery("#src").val();

    var ini = document.getElementById("src").selectionStart
    var end = document.getElementById("src").selectionEnd

    if (end > ini) {
      var multiLine = false
      for (i = ini; i < end; i++) {
        if (textAreaTxt.charAt(i)=="\n") {
          multiLine = true
        }
      }

      if (multiLine) {
        var lines = jQuery("#src").val().split('\n');
        var newLines = new Array();

        var lineIni = 0
        for (i = 0; i < lines.length; i++) {
          var f = ""
          var l = ""
          lineEnd = lineIni + lines[i].length + 1

          if (ini >= lineIni && ini < lineEnd){
            newLines[i] = "```\n"+lines[i]
          } else {
            if (end >= lineIni && end < lineEnd){
              newLines[i] = lines[i]+"\n```"
            } else {
              newLines[i] = lines[i]
            }
          }

          lineIni = lineIni + lines[i].length + 1
        }

        var output = ""
        for (i = 0; i < newLines.length; i++) {
          output = output + newLines[i]
          if (i < newLines.length-1){
            output = output +"\n"
          }
        }
        $("#src").val(output);
      } else {
        var firstPart = textAreaTxt.substring(0, i)
        var secondPart = textAreaTxt.substring(i)

        $("#src").val(textAreaTxt.substring(0, ini)+"`"+textAreaTxt.substring(ini, end)+"`"+textAreaTxt.substring(end));
      }

      update();
    }
  });

  $("#btnMD_linebreak").on('click', function() {
    var textAreaTxt = jQuery("#src").val();
    var ini = "0"+document.getElementById("src").selectionStart
    var end = "0"+document.getElementById("src").selectionEnd

    if (ini = end) {
      var extra = "\n\n"
      if (textAreaTxt.charAt(ini)=="\n") {
        extra = "\n"
      }
      $("#src").val(textAreaTxt.substring(0, ini)+"\n\n---"+extra+textAreaTxt.substring(ini));
    }
    update();
  });

  $("#btnMD_hyperlink").on('click', function() {
    var textAreaTxt = jQuery("#src").val();

    var ini = document.getElementById("src").selectionStart
    var end = document.getElementById("src").selectionEnd

    if (end > ini) {
      var multiLine = false
      for (i = ini; i < end; i++) {
        if (textAreaTxt.charAt(i)=="\n") {
          multiLine = true
        }
      }

      if (multiLine == false) {
        var firstPart = textAreaTxt.substring(0, ini)
        var selection = textAreaTxt.substring(ini, end)
        var secondPart = textAreaTxt.substring(end)
        if (isUrl(selection) == false) {
          selection = "http://"+selection
        }

        $("#src").val(firstPart+"["+selection+"]("+selection+")"+secondPart);
      }

      update();
    }
  });

  $("#btnMD_bold").on('click', function() {
    var textAreaTxt = jQuery("#src").val();

    var ini = document.getElementById("src").selectionStart
    var end = document.getElementById("src").selectionEnd
    if (end > ini) {
      var multiLine = false
      for (i = ini; i < end; i++) {
        if (textAreaTxt.charAt(i)=="\n") {
          multiLine = true
        }
      }

      if (multiLine==false) {
        $("#src").val(textAreaTxt.substring(0, ini)+"**"+textAreaTxt.substring(ini, end)+"**"+textAreaTxt.substring(end));
        update();
      }
    }
  });

  $("#btnMD_italic").on('click', function() {
    var textAreaTxt = jQuery("#src").val();

    var ini = document.getElementById("src").selectionStart
    var end = document.getElementById("src").selectionEnd
    if (end > ini) {
      var multiLine = false
      for (i = ini; i < end; i++) {
        if (textAreaTxt.charAt(i)=="\n") {
          multiLine = true
        }
      }

      if (multiLine==false) {
        $("#src").val(textAreaTxt.substring(0, ini)+"*"+textAreaTxt.substring(ini, end)+"*"+textAreaTxt.substring(end));
        update();
      }
    }
  });

  $("#btnMD_image").on('click', function() {
    var textAreaTxt = jQuery("#src").val();

    var ini = document.getElementById("src").selectionStart
    var end = document.getElementById("src").selectionEnd
    if (end > ini) {
      var multiLine = false
      for (i = ini; i < end; i++) {
        if (textAreaTxt.charAt(i)=="\n") {
          multiLine = true
        }
      }

      if (multiLine==false) {
        $("#src").val(textAreaTxt.substring(0, ini)+"![Insert alternative text]("+textAreaTxt.substring(ini, end)+")"+textAreaTxt.substring(end));
      }
    } else {
      if (ini == end) {
        $("#src").val(textAreaTxt.substring(0, ini)+"![Insert alternative text](/path/to/image)"+textAreaTxt.substring(end));
      }
    }
    update();
  });

  //From: https://jadendreamer.wordpress.com/2013/04/24/jquery-tutorial-scroll-ui-dialog-boxes-with-the-page-as-it-scrolls/
  //very smart and useful!
  $(document).scroll(function(e){

      if ($(".ui-widget-overlay")) //the dialog has popped up in modal view
      {
          //fix the overlay so it scrolls down with the page
          $(".ui-widget-overlay").css({
              position: 'fixed',
              top: '0'
          });

          //get the current popup position of the dialog box
          pos = $(".ui-dialog").position();

          //adjust the dialog box so that it scrolls as you scroll the page
          $(".ui-dialog").css({
              position: 'fixed',
              top: pos.y
          });
      }
  });

  function isUrl(s) {
     var regexp = /(ftp|http|https):\/\/.*/
     return regexp.test(s);
  }

  $("#input_title").on('input', function() {
    update_title()
  });

  $("#input_short_description").on('input', function() {
    update_short_description()
  });

  $("#src").on('input', function() {
    update()
  });

  update_title();
  update_short_description();
  update();

  var pxy = localStorage.getItem("pxy");
  if (pxy == ""){
    pxy = "center center"
  }

  $('#post-editor').dialog({
      closeOnEscape: false,
      open: function(event, ui) {
        $(".ui-dialog-titlebar-close").hide();
      },
      autoOpen: true,
      width: 640,
      height: 480,
      position: { my: "left top", at: pxy, of: window },
      modal: false,
      appendTo: "#anchor",
      buttons: [
        {
          text: "Minimize",
          click: function () {
            $("#post-editor").dialog( "option", "width", 40 );
            $("#post-editor").dialog( "option", "height", 40 );
            $("#post-editor").dialog( "position", "height", {my: "center", at: "center", of: window} );

          }
        },
        {
          text: "Maximize",
          click: function () {
            $("#post-editor").dialog( "option", "width", 640 );
            $("#post-editor").dialog( "option", "height", 480 );
          }
        },
        {
          text: "Save",
          click: function () {
            var position = $( "#post-editor" ).dialog( "option", "position" );
            localStorage.setItem("pxy", position.at);
            $("#html_content").val($("#target").html())
            $("#form").submit();
          }
        }
      ]
  });

  var objDiv = document.getElementById("body");
  objDiv.scrollTop = objDiv.scrollHeight;
  objDiv.scrollTop = 0;
});

function htmlUnescape(value){
    return String(value)
        .replace(/&quot;/g, '"')
        .replace(/&#39;/g, "'")
        .replace(/&#x2f;/g, "/")
        .replace(/&lt;/g, '<')
        .replace(/&gt;/g, '>')
        .replace(/&amp;/g, '&');
}

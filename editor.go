package main

var editor = `
<div id="anchor">
</div>
<div id="post-editor" title="Post editor" >
<div style="height:100%">
  <form id="form" method="POST" action="save" style="height:100%">
    <input type="text" name="title" id="input_title" value="{{.Title}}" placeholder="Title" style="width:100%; margin:5px 0px"/>

    <input type="text" name="short_description" id="input_short_description" value="{{.ShortDescription}}" placeholder="Short Decription" style="width:100%; margin-bottom:5px"/>
    <div style="margin-bottom:5px">
      <button id="btnMD_bold" type="button" class="btn btn-default" aria-label="Left Align" title="Bold">
        <span class="fa fa-bold fa-fw" aria-hidden="true"/>
      </button>
      <button id="btnMD_italic" type="button" class="btn btn-default" aria-label="Left Align" title="Italic">
        <span class="fa fa-italic fa-fw" aria-hidden="true"/>
      </button>
      <button id="btnMD_H1" type="button" class="btn btn-default" aria-label="Left Align" title="Header 1">
        <span class="fa fa-header fa-fw" aria-hidden="true"/>1
      </button>
      <button id="btnMD_H2" type="button" class="btn btn-default" aria-label="Left Align" title="Header 2">
        <span class="fa fa-header fa-fw" aria-hidden="true"/>2
      </button>
      <button id="btnMD_H3" type="button" class="btn btn-default" aria-label="Left Align" title="Header 3">
        <span class="fa fa-header fa-fw" aria-hidden="true"/>3
      </button>
      <button id="btnMD_image" type="button" class="btn btn-default" aria-label="Left Align" title="Image">
        <span class="fa fa-photo fa-fw" aria-hidden="true"/>
      </button>
      <button id="btnMD_hyperlink" type="button" class="btn btn-default" aria-label="Left Align" title="Hyperlink">
        <span class="fa fa-link fa-fw" aria-hidden="true"/>
      </button>


      <button id="btnMD_quote" type="button" class="btn btn-default" aria-label="Left Align" title="Quote">
        <span class="fa fa-quote-left fa-fw" aria-hidden="true"/>
      </button>
      <button id="btnMD_list" type="button" class="btn btn-default" aria-label="Left Align" title="List">
        <span class="fa fa-list-ul fa-fw" aria-hidden="true"/>
      </button>
      <button id="btnMD_orderedlist" type="button" class="btn btn-default" aria-label="Left Align" title="Ordered List">
        <span class="fa fa-list-ol fa-fw" aria-hidden="true"/>
      </button>
      <button id="btnMD_code" type="button" class="btn btn-default" aria-label="Left Align" title="Code">
        <span class="fa fa-code fa-fw" aria-hidden="true"/>
      </button>
      <button id="btnMD_linebreak" type="button" class="btn btn-default" aria-label="Left Align" title="Horizontal Line">
        <span class="fa fa-ellipsis-h fa-fw" aria-hidden="true"/>
      </button>

    </div>
    <textarea name="src_content" id="src" style="width:100%; height:90%; font-family:Consolas,Monaco,Lucida Console,Liberation Mono,DejaVu Sans Mono,Bitstream Vera Sans Mono,Courier New, monospace;">{{.Content}}</textarea>
    <input type="hidden" name="id" value="{{.ID}}"/>
    <input type="hidden" id="html_content" name="html_content" value=""/>
  </form>
</div>
</div>`

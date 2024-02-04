import { $, $htmxLoad, $load } from "../../lib/easy-dom";
import { Editor } from "../../lib/lc-wysiwyg";

$load(() => {
  let editor;

  // WYSIWYG Editor
  const setupEditor = () => {
    const editorNode = $("#lc-wysiwyg");
    if (editorNode !== null) {
      editor = new Editor(editorNode.node);
      editor.enable();
    }
  };

  document.body.addEventListener("hx:new-note", (e) => {
    console.log(e.detail.data);
  });

  window.submitNewNote = function () {
    form = $("#note-editor").node;
    formData = new FormData(form);
    sourceButton = $("button[hx-encoding='multipart/form-data']").node;

    htmx.ajax("POST", "/new", {
      values: {
        title: formData.get("title"),
        content: formData.get("content"),
      },
      source: sourceButton,
    });
  };

  $htmxLoad(setupEditor, true);
});

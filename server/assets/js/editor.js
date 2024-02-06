import { Editor } from "./lib/lc-wysiwyg";

export const setupEditor = (node, data) => {
  node.removeAttribute("x-init");
  const editor = new Editor(node);
  editor.enable();
  if (data) {
    editor.setContent(data);
  }
};

export const submitNewNote = function () {
  form = $("#note-editor").node;
  formData = new FormData(form);
  sourceButton = $("button[hx-encoding='multipart/form-data']").node;

  htmx.ajax("POST", "/new", {
    values: {
      title: formData.get("title"),
      content: formData.get("content"),
    },
    source: sourceButton,
    headers: {
      "HX-Boosted": "true",
    },
  });
};

export const updateNote = function (noteId) {
  form = $("#note-editor").node;
  formData = new FormData(form);
  sourceButton = $("button[hx-encoding='multipart/form-data']").node;

  htmx.ajax("POST", `/note/${noteId}/edit`, {
    values: {
      title: formData.get("title"),
      content: formData.get("content"),
    },
    source: sourceButton,
    headers: {
      "HX-Boosted": "true",
    },
  });
};

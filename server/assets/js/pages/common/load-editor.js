import { $, $htmxLoad, $load } from "../../lib/easy-dom";
import { Editor } from "../../lib/lc-wysiwyg";

$load(() => {
  // WYSIWYG Editor
  const setupEditor = () => {
    const editorNode = $("#lc-wysiwyg");
    if (editorNode !== null) {
      const editor = new Editor(editorNode.node);
      editor.enable();
    }
  };

  $htmxLoad(setupEditor, true);
});

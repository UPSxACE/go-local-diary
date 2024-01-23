$load(() => {
  // Toggle sidebars
  $("#showcase-header-button-sidebar").click(() => {
    $("#sidebar").toggle("sidebar-closed");
  });

  $("#showcase-header-button-info").click(() => {
    $("#main").toggle("info-sidebar-closed");
  });

  // Tabs switch
  const tabs = $all(".tab a").map((easyNode) => {
    const parent = easyNode.getParent();

    easyNode.click((e) => {
      tabs.forEach((tab) => {
        tab.toggle("tab-active", false);
      });

      parent.toggle("tab-active", true);
    });

    return parent;
  });

  $("#home-link").click(() => {
    tabs.forEach((tab) => {
      tab.toggle("tab-active", false);
    });
  });

  // Setup resizable div
  const setResizeable = () =>
    $resizeableBar($(".resize-bar").node, $("#sidebar-info-wrapper").node);

  $htmxLoad(setResizeable, true);

  // WYSIWYG Editor
  const setupEditor = () => {
    const editorNode = $("#lc-wysiwyg");
    if (editorNode !== null) {
      const editor = new Editor(editorNode.node);
      editor.setContent("<h1>Loading content test!</h1>");
      editor.enable();
    }
  };

  $htmxLoad(setupEditor, true);
});

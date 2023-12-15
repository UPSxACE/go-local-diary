$load(() => {
  $("#showcase-header-button-sidebar").click(() => {
    $("#sidebar").toggle("sidebar-closed");
  });

  $("#showcase-header-button-info").click(() => {
    $("#main").toggle("info-sidebar-closed");
  });

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

  const setResizeable = () =>
    $resizeableBar($(".resize-bar").node, $("#sidebar-info-wrapper").node);

  setResizeable();

  $htmxLoad(setResizeable);
});

// Throttling Function
window.$throttleFunction = (func, delay) => {
  // Previously called time of the function
  let prev = 0;
  return (...args) => {
    let now = new Date().getTime();
    if (now - prev > delay) {
      prev = now;
      return func(...args);
    }
  };
};

window.$resizeableBar = (nodeBar, nodeTarget) => {
  if (window.$state.infobarW) {
    nodeTarget.style.width = window.$state.infobarW;
  }
  if (!window.$state.infobarW) {
    window.$state.infobarW = `${nodeTarget.style.width}px`;
  }

  // hardware acceleration
  nodeTarget.style.transform = "translateZ(0)";
  nodeBar.style.transform = "translateZ(0)";

  let isMoving = false;
  let initialX = null;
  let initialWidth = null;

  document.addEventListener("mousedown", startedMoving);
  document.addEventListener("mousemove", duringMovement);
  document.addEventListener("mouseup", finishedMoving);

  function startedMoving(e) {
    if (e.target === nodeBar) {
      initialWidth = nodeTarget.offsetWidth;
      nodeTarget.classList.toggle("no-transition", true);
      initialX = e.clientX;
      isMoving = true;
    }
  }

  function duringMovement(e) {
    if (isMoving) {
      const difference = initialX - e.clientX;
      const newWidth = Math.max(initialWidth + difference, 10);
      window.$state.infobarW = `${newWidth}px`;
      nodeTarget.style.width = `${newWidth}px`;
    }
  }

  function finishedMoving(e) {
    isMoving = false;
    nodeTarget.classList.toggle("no-transition", false);
  }

  return;
};

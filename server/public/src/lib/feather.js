import * as feather from "feather-icons/dist/feather.min";

$load(() => {
  feather.replace();
  $htmxLoad(() => feather.replace());
});

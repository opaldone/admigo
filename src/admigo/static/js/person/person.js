;"use strict";
class PersonPage {
  constructor() {
    window.fun.ready(this.handler.bind(this));
  }

  handler() {
    Fancybox.bind("#per-ava", {});
  }
}

new PersonPage();

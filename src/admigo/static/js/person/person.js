;"use strict";

class PersonPage {
  constructor() {
    this.fun = new Funcs();
    this.fun.ready(this.handler.bind(this));
  }

  handler() {
    Fancybox.bind("#per-ava", {});
  }
}

new PersonPage();

import htmx from "htmx.org";
import Alpine from 'alpinejs'
import { BooleanCheckbox } from "./custom_elements.js";

customElements.define("boolean-checkbox", BooleanCheckbox, { extends: "input" }); //! Works fine but not I form data encoded !!

window.Alpine = Alpine;
Alpine.start()
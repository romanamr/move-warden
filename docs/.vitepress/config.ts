import { defineConfig } from "vitepress";

export default defineConfig({
  title: "MoveWarden",
  description: "Documentación de uso / Usage documentation",
  cleanUrls: true,
  themeConfig: {
    search: {
      provider: "local",
    },
  },
  locales: {
    root: {
      label: "Root",
      lang: "es",
      link: "/",
      themeConfig: {
        nav: [
          { text: "Español", link: "/es/" },
          { text: "English", link: "/en/" },
        ],
      },
    },
    es: {
      label: "Español",
      lang: "es-ES",
      link: "/es/",
      themeConfig: {
        nav: [
          { text: "Inicio", link: "/es/" },
          { text: "Guía", link: "/es/guide/quick-start" },
          { text: "Referencia", link: "/es/reference/rules-format" },
          { text: "English", link: "/en/" },
        ],
        sidebar: [
          {
            text: "Guía de uso",
            items: [
              { text: "Inicio rápido", link: "/es/guide/quick-start" },
              { text: "Transformaciones", link: "/es/guide/transformations" },
              { text: "Filtros", link: "/es/guide/filters" },
            ],
          },
          {
            text: "Referencia",
            items: [
              {
                text: "Formato de rules.json",
                link: "/es/reference/rules-format",
              },
              {
                text: "Variables y placeholders",
                link: "/es/reference/variables",
              },
              { text: "Ejemplos", link: "/es/reference/examples" },
            ],
          },
        ],
      },
    },
    en: {
      label: "English",
      lang: "en-US",
      link: "/en/",
      themeConfig: {
        nav: [
          { text: "Home", link: "/en/" },
          { text: "Guide", link: "/en/guide/quick-start" },
          { text: "Reference", link: "/en/reference/rules-format" },
          { text: "Español", link: "/es/" },
        ],
        sidebar: [
          {
            text: "Usage Guide",
            items: [
              { text: "Quick Start", link: "/en/guide/quick-start" },
              { text: "Transformations", link: "/en/guide/transformations" },
              { text: "Filters", link: "/en/guide/filters" },
            ],
          },
          {
            text: "Reference",
            items: [
              { text: "rules.json format", link: "/en/reference/rules-format" },
              {
                text: "Variables and placeholders",
                link: "/en/reference/variables",
              },
              { text: "Examples", link: "/en/reference/examples" },
            ],
          },
        ],
      },
    },
  },
  base: "/move-warden/",
});

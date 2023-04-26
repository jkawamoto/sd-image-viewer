import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import {MantineProvider} from "@mantine/core";

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <MantineProvider theme={{colorScheme: 'dark'}} withGlobalStyles withNormalizeCSS>
      <App/>
    </MantineProvider>
  </React.StrictMode>,
)

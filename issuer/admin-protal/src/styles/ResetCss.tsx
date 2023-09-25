import { FC } from "react"
import { Global, css } from "@emotion/react"
import normalizeCss from "./normalize"

export const ResetCss: FC = () => {
  return (
    <Global
      styles={css`
        ${normalizeCss.styles}
        * {
          padding: 0;
          margin: 0;
          box-sizing: border-box;
        }
        html,
        body {
          min-height: 100vh;
          min-width: 1300px;
        }
      `}
    />
  )
}

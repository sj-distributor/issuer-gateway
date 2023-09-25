import { CSSProperties } from "react"
import { CSSObject } from "@emotion/react"
import { singleLineEllipsis } from "@/styles/base"

export const domainCell: CSSObject = {
  ...singleLineEllipsis,
  width: "fit-content",
  maxWidth: "100%",
}

export const helpTextStyles: CSSProperties = {
  marginTop: 0,
}

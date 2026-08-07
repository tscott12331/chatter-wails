import { ComponentPropsWithoutRef } from "react"

const EmoteIcon = (props: ComponentPropsWithoutRef<'svg'>) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    shapeRendering="geometricPrecision"
    textRendering="geometricPrecision"
    viewBox="0 0 300 300"
    {...props}
  >
    <circle r={30} fill="#f4effa" transform="matrix(5 0 0 5 150 150)" />
    <ellipse
      fill="#342a40"
      rx={30}
      ry={30.26}
      transform="matrix(1 0 0 1.23927 90 97.5)"
    />
    <ellipse
      fill="#342a40"
      rx={30}
      ry={30.26}
      transform="matrix(1 0 0 1.23927 210 97.5)"
    />
    <g mask="url(#a)" transform="rotate(180 150.04 200.5)">
      <circle
        r={30}
        fill="#342a40"
        transform="matrix(3.00133 0 0 2.3 150.04 208)"
      />
      <mask
        id="a"
        width="400%"
        height="400%"
        x="-150%"
        y="-150%"
        mask-type="luminance"
      >
        <path
          fill="#fff"
          strokeWidth={0}
          d="M240.04 208c0 10.493-40.294 19-90 19s-90-8.507-90-19H60v-69h180v68.428c.027.19.04.38.04.572Z"
        />
      </mask>
    </g>
    <g mask="url(#b)" transform="matrix(-.73304 0 0 -.40057 259.982 291.106)">
      <mask
        id="b"
        width="400%"
        height="400%"
        x="-150%"
        y="-150%"
        mask-type="luminance"
      >
        <path
          fill="#fff"
          strokeWidth={0}
          d="M240.04 208c0 10.493-40.294 19-90 19s-90-8.507-90-19H60v-69h180v68.428c.027.19.04.38.04.572Z"
        />
      </mask>
    </g>
  </svg>
)
export default EmoteIcon

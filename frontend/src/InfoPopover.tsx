const POPOVER_ID = 'info-popover'

// Uses the native Popover API: the browser opens and closes it (including Esc and
// clicking outside) and manages focus order, so no React state is needed.
export function InfoPopover() {
  return (
    <>
      <button type="button" className="info-button" popoverTarget={POPOVER_ID} aria-label="How to use">
        i
      </button>

      <div id={POPOVER_ID} className="info-popover" popover="auto">
        <h2>How to use</h2>
        <ol>
          <li>Enter the first number.</li>
          <li>Pick an operation: + − × ÷ xʸ √ %. With %, 10 and 200 give 10% of 200 = 20.</li>
          <li>Enter the second number (not needed for √) and press Calculate or Enter.</li>
        </ol>

        <h2 id="info-valid">Valid numbers</h2>
        <ul aria-labelledby="info-valid">
          <li>
            Whole numbers: <code>42</code>, <code>-7</code>, <code>+3</code>
          </li>
          <li>
            Decimals: <code>2.5</code>, <code>.5</code>, <code>5.</code>
          </li>
          <li>
            Scientific notation: <code>1e5</code>, <code>2.5e-3</code>
          </li>
        </ul>

        <h2 id="info-invalid">Not accepted</h2>
        <ul aria-labelledby="info-invalid">
          <li>
            Commas: <code>1,000</code>, <code>1,5</code>
          </li>
          <li>
            Text, hex or special values: <code>abc</code>, <code>0x10</code>, <code>Infinity</code>
          </li>
          <li>
            Numbers too large to represent: <code>1e400</code>
          </li>
        </ul>

        <button type="button" className="info-close" popoverTarget={POPOVER_ID} popoverTargetAction="hide">
          Close
        </button>
      </div>
    </>
  )
}

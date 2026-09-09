import "../styles/Components.styles.css"

export default function FloatingInput({
  label,
  id,
  value,
  onChange,
  type = "text",
  ...props
}) {
  return (
    <div className="floating-field">
      <input
        id={id}
        type={type}
        value={value}
        onChange={onChange}
        placeholder=" "
        {...props}
      />

      <label htmlFor={id}>{label}</label>
    </div>
  );
}

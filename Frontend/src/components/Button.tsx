interface ButtonProps {
  children: React.ReactNode;
  onClick?: () => void;
  disabled?: boolean;
  isLoading?: boolean;
  className?: string
}

export const Button: React.FC<ButtonProps> = ({
  children,
  onClick,
  disabled,
  isLoading,
  className,
}) => {
  return (
    <button
      onClick={onClick}
      disabled={disabled || isLoading}
      className={` cursor-pointer whitespace-nowrap rounded-full ${className}`}
    >
      {isLoading ? "Loading..." : children}
    </button>
  );
};



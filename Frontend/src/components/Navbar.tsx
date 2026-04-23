import Logo from "./Logo";
import Search from "./Search";
import { Button } from "./Button";
function Navbar() {
  return (
    <nav className="sticky top-0 w-full bg-white shadow-sm ">
      <div className="flex justify-center align-center gap-6 px-4 py-4">
        <Logo />
        <Search />

        <Button className=" bg-[#34967C] text-white">Sign In</Button>
        <Button className="border border-[#ECECEC]">Sign Up</Button>
      </div>
    </nav>
  );
}

export default Navbar;

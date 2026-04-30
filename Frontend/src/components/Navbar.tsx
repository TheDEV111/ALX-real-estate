import Logo from "./Logo";
import Search from "./Search";
import { Button } from "./Button";

function Navbar() {
  return (
    <nav className="sticky top-0 w-full bg-white shadow-sm font-[Quicksand] ">
      <div className="flex items-center gap-8 px-8 py-4">
        <Logo />
        <div className="flex-1 flex justify-center">
          <Search />
        </div>

        <div className="flex items-center gap-4 shrink-0">
          <Button className="bg-[#34967C] text-white rounded-full px-6 hover:bg-[#2a7a65] transition-colors">
            Sign In
          </Button>
          <Button className="border border-[#ECECEC] px-6 hover:bg-gray-50 transition-colors">
            Sign Up
          </Button>
        </div>
      </div>
    </nav>
  );
}

export default Navbar;

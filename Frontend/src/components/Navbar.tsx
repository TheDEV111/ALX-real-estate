import Logo from "./logo";
import Search from "./Search";
function Navbar() {
  return (
    <nav className="fixed w-full bg-white shadow-sm ">
      <div className="flex justify-center align-center gap-6 max-w-[2520px] mx-auto xl:px-20 md:px-10 sm:px-2 px-4 py-4">
        <Logo />
        <Search />
      </div>
    </nav>
  );
}

export default Navbar;

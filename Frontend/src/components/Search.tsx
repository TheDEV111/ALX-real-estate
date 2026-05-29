import { useState } from "react";
import DatePicker from "react-datepicker";
import "react-datepicker/dist/react-datepicker.css";
import "../styles/datepicker.css";
import { BiSearch } from "react-icons/bi";

function Search() {
  const [location, setLocation] = useState<string>("");
  const [checkIn, setCheckIn] = useState<Date | null>(null);
  const [checkOut, setCheckOut] = useState<Date | null>(null);
  const [people, setPeople] = useState<number>();

  const handlePeopleChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
    const value = parseInt(e.target.value);
    if (value >= 1) {
      setPeople(value);
    }
  };

  const handleSearch = (): void => {
    console.log({ location, checkIn, checkOut, people });
  };

  return (
    <div className="flex gap-10 items-center  ">
      <div className="border border-[#F6F6F6]  py-1 rounded-full shadow-sm hover:shadow-md transition cursor-pointer flex items-center">
        {/* Mobile Search */}
        <div className="md:hidden text-sm  flex flex-col gap-3  ml-8">
          <p className="font-semibold whitespace-nowrap">Where to</p>
          <div className="flex items-center">
            <div className="flex items-center pr-2 ">
              <input
                type="text"
                placeholder="Location"
                value={location}
                onChange={(e) => setLocation(e.target.value)}
                className="focus:outline-none placeholder:text-[12px] bg-transparent text-[14px] field-sizing-content pr-3"
              />
              <div className="w-[5px] h-[5px] bg-gray-500 rounded-full  "></div>
            </div>
            <div className="flex items-center gap-1">
              <DatePicker
                selected={checkIn}
                onChange={(date: Date | null) => setCheckIn(date)}
                placeholderText="Date"
                className="focus: placeholder:text-[12px] bg-transparent text-[14px] cursor-pointer field-sizing-content "
                minDate={new Date()}
              />
              <div className="w-[5px] h-[5px] bg-gray-500 rounded-full  "></div>
            </div>
            <input
              type="number"
              placeholder="Add guest"
              value={people}
              onChange={handlePeopleChange}
              min="1"
              className=" placeholder:text-[12px] bg-transparent text-[14px] w-20 field-sizing-content pr-3"
            />
          </div>
        </div>

        {/* Desktop Search */}
        <div className="hidden md:flex flex-1">
          <div className="text-sm px-6 flex flex-col border-r border-gray-200">
            <label className="font-semibold">Location</label>
            <input
              type="text"
              placeholder="Search for destination"
              value={location}
              onChange={(e) => setLocation(e.target.value)}
              className="focus:outline-none placeholder:text-[12px] bg-transparent text-[14px]"
            />
          </div>

          <div className="text-sm px-6 flex flex-col border-r border-gray-200">
            <label className="font-semibold">Check in</label>
            <DatePicker
              selected={checkIn}
              onChange={(date: Date | null) => setCheckIn(date)}
              placeholderText="Add date"
              className="focus:outline-none placeholder:text-[12px] bg-transparent text-[14px] cursor-pointer"
              minDate={new Date()}
            />
          </div>

          <div className="text-sm px-6 flex flex-col border-r border-gray-200">
            <label className="font-semibold">Check out</label>
            <DatePicker
              selected={checkOut}
              onChange={(date: Date | null) => setCheckOut(date)}
              placeholderText="Add date"
              className="focus:outline-none placeholder:text-[12px] bg-transparent text-[14px] cursor-pointer"
              minDate={checkIn || new Date()}
            />
          </div>

          <div className="text-sm px-6 flex flex-col">
            <label className="font-semibold">People</label>
            <input
              type="number"
              placeholder="Add guests"
              value={people}
              onChange={handlePeopleChange}
              min="1"
              className="focus:outline-none placeholder:text-[12px] bg-transparent text-[14px] w-20"
            />
          </div>
        </div>

        {/* Search Button */}
        <button
          onClick={handleSearch}
          className="p-3 bg-[#FFA800] rounded-full text-white hover:bg-[#FF9500] transition cursor-pointer ml-2 mr-2 flex-shrink-0"
        >
          <BiSearch size={18} />
        </button>
      </div>
      <img
        src="/icons/Vector.svg"
        alt="ALX Logo"
        className="w-8 h-8 md:hidden bg-[#34967C] rounded-full"
      />
    </div>
  );
}

export default Search;

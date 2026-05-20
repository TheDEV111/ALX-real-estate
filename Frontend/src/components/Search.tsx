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
      setPeople(value); // Store as number
    }
  };
  
 const handleSearch = (): void => {
   // Handle search logic here
   console.log({ location, checkIn, checkOut, people });
 };

  return (
    <div className="border-[#F6F6F6] w-full md:w-auto py-2 rounded-full shadow-sm hover:shadow-md transition cursor-pointer flex ">
      <div className="text-sm px-6  flex flex-col">
        <label className="font-semibold">Location</label>
        <input
          type="text"
          placeholder="Search for destination"
          value={location}
          onChange={(e) => setLocation(e.target.value)}
          className="focus:outline-none placeholder:text-[12px] bg-transparent text-[14px]"
        />
      </div>

      <div className="text-sm px-6  flex flex-col">
        <label className="font-semibold">Check in</label>
        <DatePicker
          selected={checkIn}
          onChange={(date: Date | null) => setCheckIn(date)}
          placeholderText="Add date"
          className="focus:outline-none placeholder:text-[12px] bg-transparent text-[14px] cursor-pointer"
          minDate={new Date()}
        />
      </div>

      <div className="text-sm px-6  flex flex-col">
        <label className="font-semibold">Check out</label>
        <DatePicker
          selected={checkOut}
          onChange={(date: Date | null) => setCheckOut(date)}
          placeholderText="Add date"
          className="focus:outline-none placeholder:text-[12px] bg-transparent text-[14px] cursor-pointer"
          minDate={checkIn || new Date()}
        />
      </div>

      <div className="text-sm px-10   flex flex-col">
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

      <button
        onClick={handleSearch}
        className="p-3 bg-[#FFA800] rounded-full text-white hover:bg-[#FF9500] transition cursor-pointer ml-2 mr-2"
      >
        <BiSearch size={18} />
      </button>
    </div>
  );
}

export default Search;

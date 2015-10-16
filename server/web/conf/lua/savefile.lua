package.path = '/usr/local/share/lua/5.1/?.lua;/usr/local/openresty/lualib/resty/?.lua;'  
package.cpath = '/usr/local/lib/lua/5.1/?.so;'  


do
   -- declare local variables
   --// exportstring( string )
   --// returns a "Lua" portable version of the string
   local function exportstring( s )
      return string.format("%q", s)
   end

   --// The Save Function
   function table.save(  tbl,filename )
      local charS,charE = "   ","\n"
      local file,err = io.open( filename, "wb" )
      if err then return err end

      -- initiate variables for save procedure
      local tables,lookup = { tbl },{ [tbl] = 1 }
      file:write( "return {"..charE )

      for idx,t in ipairs( tables ) do
         file:write( "-- Table: {"..idx.."}"..charE )
         file:write( "{"..charE )
         local thandled = {}

         for i,v in ipairs( t ) do
            thandled[i] = true
            local stype = type( v )
            -- only handle value
            if stype == "table" then
               if not lookup[v] then
                  table.insert( tables, v )
                  lookup[v] = #tables
               end
               file:write( charS.."{"..lookup[v].."},"..charE )
            elseif stype == "string" then
               file:write(  charS..exportstring( v )..","..charE )
            elseif stype == "number" then
               file:write(  charS..tostring( v )..","..charE )
            end
         end

         for i,v in pairs( t ) do
            -- escape handled values
            if (not thandled[i]) then
            
               local str = ""
               local stype = type( i )
               -- handle index
               if stype == "table" then
                  if not lookup[i] then
                     table.insert( tables,i )
                     lookup[i] = #tables
                  end
                  str = charS.."[{"..lookup[i].."}]="
               elseif stype == "string" then
                  str = charS.."["..exportstring( i ).."]="
               elseif stype == "number" then
                  str = charS.."["..tostring( i ).."]="
               end
            
               if str ~= "" then
                  stype = type( v )
                  -- handle value
                  if stype == "table" then
                     if not lookup[v] then
                        table.insert( tables,v )
                        lookup[v] = #tables
                     end
                     file:write( str.."{"..lookup[v].."},"..charE )
                  elseif stype == "string" then
                     file:write( str..exportstring( v )..","..charE )
                  elseif stype == "number" then
                     file:write( str..tostring( v )..","..charE )
                  end
               end
            end
         end
         file:write( "},"..charE )
      end
      file:write( "}" )
      file:close()
   end
   
   --// The Load Function
   function table.load( sfile )
      local ftables,err = loadfile( sfile )
      if err then return _,err end
      local tables = ftables()
      for idx = 1,#tables do
         local tolinki = {}
         for i,v in pairs( tables[idx] ) do
            if type( v ) == "table" then
               tables[idx][i] = tables[v[1]]
            end
            if type( i ) == "table" and tables[i[1]] then
               table.insert( tolinki,{ i,tables[i[1]] } )
            end
         end
         -- link indices
         for _,v in ipairs( tolinki ) do
            tables[idx][v[2]],tables[idx][v[1]] =  tables[idx][v[1]],nil
         end
      end
      return tables[1]
   end
-- close do
end


--测试方法：在命令行执行curl -F file=@0_0_1.zip http://127.0.0.1:8080/uploadfile
local upload = require "upload"  
  
local chunk_size = 4096  
local form = upload:new(chunk_size)  
local file  
local filelen=0  
form:set_timeout(0) -- 1 sec  
local filename  
  
function get_filename(res)  
    local filename = ngx.re.match(res,'(.+)filename="(.+)"(.*)')  
    if filename then   
        return filename[2]  
    end  
end  

function getFileName(str)
    local idx = str:match(".+()%.%w+$")
    if(idx) then
        return str:sub(1, idx-1)
    else
        return str
    end
end

function IsInTable(value, tbl)
    for k,v in ipairs(tbl) do
        if v == value then
        return true;
        end
    end
    return false;
end
--工作目录
local wd = "/home/jackie/work/winner/server/web"
--文件上传写入目录
local osfilepath = wd .. "/data/"
local i=0  
while true do  
    local typ, res, err = form:read()  
    if not typ then  
        ngx.say("failed to read: ", err)  
        return  
    end  
    if typ == "header" then  
        if res[1] ~= "Content-Type" then  
            filename = get_filename(res[2])  
            if filename then  
                i=i+1  
                filepath = osfilepath  .. filename  
                file = io.open(filepath,"w+")  
                if not file then  
                    ngx.say("failed to open file ")  
                    return  
                end  
            else  
            end  
        end  
    elseif typ == "body" then  
        if file then  
            filelen= filelen + tonumber(string.len(res))      
            file:write(res)  
        else  
        end  
    elseif typ == "part_end" then  
        if file then  
            file:close()  
            file = nil
            ngx.say("file upload success") 

            --save version
            local version_list_path = wd .. "/conf/lua/version_list.lua"
            local version_str = getFileName(string.gsub(filename, "_", "."))
            local t,err = table.load(version_list_path)
            assert( err == nil ) 
            if not IsInTable(version_str,t) then
                table.insert(t, version_str)
            end
            table.save(t, version_list_path)
        end
    elseif typ == "eof" then  
        break  
    else  
    end  
end  
if i==0 then  
    ngx.say("please upload at least one file!")  
    return  
end  
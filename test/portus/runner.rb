# Copyright (C) 2017-2018 Miquel Sabaté Solà <msabate@suse.com>
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

#
# First of all, clean up the environment.
#

ActiveRecord::Base.establish_connection
ActiveRecord::Base.connection.execute('SET FOREIGN_KEY_CHECKS = 0')
ActiveRecord::Base.connection.tables.each do |table|
  next if table == 'schema_migrations'

  ActiveRecord::Base.connection.execute("TRUNCATE #{table}")
end
ActiveRecord::Base.connection.execute('SET FOREIGN_KEY_CHECKS = 1')

unless ARGV.first == 'nopopulate'
  #
  # This registry is fake, but we need to create something :^)
  #

  Registry.create!(
    name:     'registry',
    hostname: 'registry:5000',
    use_ssl:  false
  )

  #
  # Create main user and an application token so it can be used by the test
  # suite.
  #

  user = User.create!(
    username: 'admin',
    password: '12341234',
    email:    'admin@example.local',
    admin:    true
  )

  _, plain_token = ApplicationToken.create_token(
    current_user: user,
    params:       { application: 'app' }
  )

  #
  # Create some fake repo & tag for testing purposes.
  #

  repo = Repository.create!(namespace: user.namespace, name: 'fake')
  Tag.create!(name: 'tag1', repository: repo, author: user, digest: 'digest',
              image_id: 'imageid')
end

#
# Write contents to config file.
#

config = "export PORTUSCTL_API_USER=admin\n" \
         "export PORTUSCTL_API_TOKEN=#{plain_token}\n" \
         "export PORTUSCTL_API_SERVER=http://localhost:3000\n" \
         "export PORTUSCTL=#{ENV['PORTUSCTL']}\n"

File.write('/srv/Portus/tmp/config.sh', config)
File.chmod(0o777, '/srv/Portus/tmp/config.sh')

#
# Clean coverage directory for Docker executions.
#

require 'fileutils'
FileUtils.mkdir_p('/srv/Portus/tmp/coverage', mode: 0o777)

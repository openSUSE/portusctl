#
# First of all, clean up the environment.
#

require 'database_cleaner'
DatabaseCleaner.clean_with :truncation

#
# This registry is fake, but we need to create something :^)
#

Registry.create!(
  name:     'registry',
  hostname: 'registry:5000',
  use_ssl:  false
)

#
# Create main user and an application token so it can be used by the test suite.
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

#
# Write contents to config file.
#

config = "export PORTUSCTL_API_USER=admin\n" \
         "export PORTUSCTL_API_TOKEN=#{plain_token}\n" \
         "export PORTUSCTL_API_SERVER=http://localhost:3000\n" \
         "export PORTUSCTL=#{ENV['PORTUSCTL']}\n"

File.write('/srv/Portus/tmp/config.sh', config)
File.chmod(0o777, '/srv/Portus/tmp/config.sh')
